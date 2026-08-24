package install

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Xoje-Tech/xoje-environment/internal/config"
	"github.com/Xoje-Tech/xoje-environment/internal/tui/styles"
)

type ToolType string

const (
	GoModule      ToolType = "go"
	GitHubRelease ToolType = "github"
)

type Tool struct {
	Name        string   `json:"name"`
	Type        ToolType `json:"type"`
	Module      string   `json:"module,omitempty"` // For Go
	Repo        string   `json:"repo,omitempty"`   // For GitHub
	BinaryName  string   `json:"binary_name"`      // Expected binary name in PATH
	Description string   `json:"description"`
}

var Registry = []Tool{
	{
		Name:        "gentle-ai",
		Type:        GoModule,
		Module:      "github.com/gentleman-programming/gentle-ai/v2/cmd/gentle-ai@latest",
		BinaryName:  "gentle-ai",
		Description: "The core AI ecosystem manager",
	},
	{
		Name:        "dev-tracker",
		Type:        GitHubRelease,
		Repo:        "Xoje-Tech/dev-tracker",
		BinaryName:  "dt",
		Description: "Project management and Kanban CLI",
	},
}

// ToolStatus represents the current state of a tool
type ToolStatus struct {
	Tool      Tool
	Installed bool
	Path      string
	Version   string
}

// DefaultBinDir returns the canonical user binary directory (~/.local/bin).
func DefaultBinDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "bin")
}

func IsInstalled(binaryName string) (string, bool) {
	path, err := exec.LookPath(binaryName)
	return path, err == nil
}

// GetStatus checks the current status of a tool
func GetStatus(t Tool) ToolStatus {
	status := ToolStatus{Tool: t}
	path, err := exec.LookPath(t.BinaryName)
	if err == nil {
		status.Installed = true
		status.Path = path
		status.Version = GetVersion(t.BinaryName)
	}
	return status
}

// GetVersion attempts to get the version of the tool.
func GetVersion(name string) string {
	cmd := exec.Command(name, "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		cmd = exec.Command(name, "--version")
		output, err = cmd.CombinedOutput()
		if err != nil {
			return "unknown"
		}
	}

	v := string(output)
	// Handle complex version strings (e.g., "gentle-ai 2.2.4")
	if strings.HasPrefix(v, name+" ") {
		v = strings.TrimPrefix(v, name+" ")
	}

	if strings.Contains(v, "{") { // Handle JSON output from dt
		var data struct {
			Version string `json:"version"`
		}
		if err := json.Unmarshal(output, &data); err == nil {
			return data.Version
		}
	}

	return strings.TrimSpace(v)
}

func InstallTool(t Tool, customBinDir string, isUpdate bool) error {
	targetDir := customBinDir
	if targetDir == "" {
		targetDir = DefaultBinDir()
	}

	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("failed to resolve target path: %w", err)
	}

	if err := os.MkdirAll(absTarget, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", absTarget, err)
	}

	fmt.Printf("%s %s...\n", styles.IrisStyle.Render(ifThen(isUpdate, "Checking/Updating", "Installing")), styles.IrisStyle.Bold(true).Render(t.Name))

	switch t.Type {
	case GoModule:
		return installGoTool(t, absTarget, isUpdate)
	case GitHubRelease:
		return installGitHubTool(t, absTarget, isUpdate)
	default:
		return fmt.Errorf("unsupported tool type: %s", t.Type)
	}
}

func installGoTool(t Tool, targetDir string, isUpdate bool) error {
	if isUpdate {
		remoteVer, err := getLatestGoModuleVersion(t.Module)
		if err == nil {
			currentVer := GetVersion(t.BinaryName)
			cleanRemote := strings.TrimPrefix(remoteVer, "v")
			cleanCurrent := strings.TrimPrefix(currentVer, "v")

			if cleanCurrent != "unknown" && cleanCurrent == cleanRemote {
				fmt.Printf("  %s %s is already up to date (%s).\n", styles.SuccessStyle.Render("✓"), t.Name, remoteVer)
				return nil
			}
			fmt.Printf("  %s Found newer version %s (current: %s)\n", styles.MutedStyle.Render("↪"), styles.IrisStyle.Render(remoteVer), cleanCurrent)
		}
	}

	cmd := exec.Command("go", "install", t.Module)
	cmd.Env = append(os.Environ(), "GOBIN="+targetDir)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go install failed: %w (output: %s)", err, string(output))
	}

	if path, ok := exec.LookPath(t.BinaryName); ok == nil {
		if t.Name == "gentle-ai" && isUpdate {
			cascadeGentleAI(path)
		} else {
			fmt.Printf("  %s %s is ready.\n", styles.SuccessStyle.Render("✓"), t.Name)
		}
	}

	return nil
}

func getLatestGoModuleVersion(modulePath string) (string, error) {
	// Extract module base: remove @latest and anything after cmd/
	base := strings.Split(modulePath, "@")[0]
	if idx := strings.Index(base, "/cmd/"); idx != -1 {
		base = base[:idx]
	}

	cmd := exec.Command("go", "list", "-m", "-json", base+"@latest")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	var data struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(output, &data); err != nil {
		return "", err
	}
	return data.Version, nil
}

func installGitHubTool(t Tool, targetDir string, isUpdate bool) error {
	assetURL, tagName, err := getLatestReleaseAsset(t.Repo, t.BinaryName)
	if err != nil {
		return err
	}

	currentVersion := GetVersion(t.BinaryName)
	cleanTarget := strings.TrimPrefix(tagName, "v")
	cleanCurrent := strings.TrimPrefix(currentVersion, "v")

	isMatch := cleanCurrent != "unknown" && (cleanCurrent == cleanTarget || strings.Contains(tagName, cleanCurrent))

	if isUpdate && isMatch {
		fmt.Printf("  %s %s is already up to date (%s).\n", styles.SuccessStyle.Render("✓"), t.Name, tagName)
		return nil
	}

	fmt.Printf("  %s Found version %s\n", styles.MutedStyle.Render("↪"), styles.IrisStyle.Render(tagName))

	targetPath := filepath.Join(targetDir, t.BinaryName)
	if err := downloadFile(assetURL, targetPath); err != nil {
		return err
	}

	if err := os.Chmod(targetPath, 0755); err != nil {
		return err
	}

	fmt.Printf("  %s Successfully %s %s to %s\n", styles.SuccessStyle.Render("✓"), ifThen(isUpdate, "updated", "installed"), t.Name, tagName)
	return nil
}

func getLatestReleaseAsset(repo, binaryName string) (string, string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases", repo)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("github api returned %d", resp.StatusCode)
	}

	var releases []struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return "", "", err
	}

	searchPattern := fmt.Sprintf("-%s-%s", runtime.GOOS, ifThen(runtime.GOARCH == "amd64", "x64", runtime.GOARCH))

	for _, release := range releases {
		for _, asset := range release.Assets {
			if strings.Contains(asset.Name, searchPattern) {
				return asset.BrowserDownloadURL, release.TagName, nil
			}
		}
	}

	return "", "", fmt.Errorf("no asset found for %s/%s in recent releases", runtime.GOOS, runtime.GOARCH)
}

func downloadFile(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func cascadeGentleAI(path string) {
	fmt.Printf("%s Triggering native gentle-ai ecosystem upgrade...\n", styles.MutedStyle.Render("↪"))
	exec.Command(path, "upgrade").Run()
	fmt.Printf("%s Ecosystem tools upgraded.\n", styles.SuccessStyle.Render("✓"))

	fmt.Printf("%s Syncing ecosystem configs and skills...\n", styles.MutedStyle.Render("↪"))
	exec.Command(path, "sync").Run()
	fmt.Printf("%s Ecosystem synchronized.\n", styles.SuccessStyle.Render("✓"))
}

func Install(name string, cfg *config.Config) error {
	for _, t := range Registry {
		if t.Name == name {
			if err := InstallTool(t, "", false); err != nil {
				return err
			}
			registerTool(name, cfg)
			return nil
		}
	}
	return fmt.Errorf("tool not found: %s", name)
}

func Update(name string, cfg *config.Config) error {
	if name == "" {
		return UpdateAll(cfg)
	}
	for _, t := range Registry {
		if t.Name == name {
			if err := InstallTool(t, "", true); err != nil {
				return err
			}
			registerTool(name, cfg)
			return nil
		}
	}
	return fmt.Errorf("tool not found: %s", name)
}

func registerTool(name string, cfg *config.Config) {
	for _, installed := range cfg.InstalledTools {
		if installed == name {
			return
		}
	}
	cfg.InstalledTools = append(cfg.InstalledTools, name)
	if err := cfg.Save(); err != nil {
		fmt.Printf("⚠️  Warning: failed to update config registry: %v\n", err)
	}
}

func UpdateAll(cfg *config.Config) error {
	fmt.Printf("%s\n", styles.TitleStyle.Render("--- 📦 Updating Registered Tools ---"))
	var failures []string
	for _, name := range cfg.InstalledTools {
		if err := Update(name, cfg); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", name, err))
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("some updates failed:\n%s", strings.Join(failures, "\n"))
	}
	return nil
}

func ifThen(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}
