package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Xoje-Tech/xoje-environment/internal/tui/styles"
)

var whitelist = map[string]string{
	"gentle-ai": "github.com/gentleman-programming/gentle-ai/v2/cmd/gentle-ai@latest",
}

// DefaultBinDir returns the canonical user binary directory (~/.local/bin).
func DefaultBinDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "bin")
}

// IsInstalled returns the absolute path of the tool if found in PATH.
func IsInstalled(name string) (string, bool) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", false
	}
	return path, true
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
	return strings.TrimSpace(string(output))
}

func InstallTool(name string, customBinDir string, isUpdate bool) error {
	module, ok := whitelist[name]
	if !ok {
		return fmt.Errorf("unknown tool: %s", name)
	}

	path, exists := IsInstalled(name)
	if exists && !isUpdate {
		fmt.Printf("%s %s is already installed at: %s\n", styles.SuccessStyle.Render("✓"), name, styles.MutedStyle.Render(path))
		return nil 
	}

	if isUpdate {
		fmt.Printf("Updating %s from %s...\n", styles.IrisStyle.Render(name), styles.MutedStyle.Render(module))
	} else {
		fmt.Printf("Installing %s from %s...\n", styles.IrisStyle.Render(name), styles.MutedStyle.Render(module))
	}

	// Determine the target directory
	targetDir := customBinDir
	if targetDir == "" {
		targetDir = DefaultBinDir()
	}

	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("failed to resolve target path: %w", err)
	}

	// Ensure the directory exists
	if err := os.MkdirAll(absTarget, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", absTarget, err)
	}

	cmd := exec.Command("go", "install", module)
	cmd.Env = append(os.Environ(), "GOBIN="+absTarget)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go install failed: %w (output: %s)", err, string(output))
	}

	// Post-install verification
	if newPath, verified := IsInstalled(name); verified {
		fmt.Printf("%s %s verified at %s: %s\n", styles.SuccessStyle.Render("✓"), name, styles.MutedStyle.Render(newPath), GetVersion(name))
	}

	return nil
}
