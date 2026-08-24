package doctor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Xoje-Tech/xoje-environment/internal/install"
)

// GoVersionCheck verifies the installed Go version.
type GoVersionCheck struct {
	Required string
}

func (c *GoVersionCheck) ID() string   { return "go-version" }
func (c *GoVersionCheck) Name() string { return "Go Version" }
func (c *GoVersionCheck) Run() Result {
	res := Result{Name: c.Name(), Status: StatusPass}
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		res.Status = StatusFail
		res.Detail = "Go is not installed or not in PATH."
		res.Remedy = "Install Go 1.22+ from https://go.dev/dl/"
		return res
	}

	version := string(out)
	res.Detail = strings.TrimSpace(version)
	if !strings.Contains(version, "go"+c.Required) && !strings.Contains(version, "go1.26") {
		res.Status = StatusWarn
		res.Remedy = "Upgrade to Go " + c.Required + "+ for best compatibility."
	}

	return res
}

// ExecutableCheck verifies that a required runtime executable is available.
type ExecutableCheck struct {
	Binary string
	Label  string
	Remedy string
}

func (c *ExecutableCheck) ID() string   { return c.Binary + "-runtime" }
func (c *ExecutableCheck) Name() string { return c.Label + " Runtime" }
func (c *ExecutableCheck) Run() Result {
	result := Result{Name: c.Name(), Status: StatusPass}
	path, err := exec.LookPath(c.Binary)
	if err != nil {
		result.Status = StatusFail
		result.Detail = c.Label + " is not installed or not in PATH."
		result.Remedy = c.Remedy
		return result
	}
	result.Detail = path
	return result
}

// EnvPathCheck verifies that xoje's canonical binary directory is in PATH.
type EnvPathCheck struct{}

func (c *EnvPathCheck) ID() string   { return "env-path" }
func (c *EnvPathCheck) Name() string { return "System PATH Configuration" }
func (c *EnvPathCheck) Run() Result {
	res := Result{Name: c.Name(), Status: StatusPass}
	binDir := install.DefaultBinDir()
	canonicalBin := canonicalPath(binDir)

	found := false
	for _, path := range filepath.SplitList(os.Getenv("PATH")) {
		if path != "" && canonicalPath(path) == canonicalBin {
			found = true
			break
		}
	}

	if !found {
		res.Status = StatusFail
		res.Detail = fmt.Sprintf("xoje binary directory (%s) is NOT in your PATH.", binDir)
		res.Remedy = fmt.Sprintf("Add 'export PATH=$PATH:%s' to your shell profile.", binDir)
		return res
	}

	res.Detail = fmt.Sprintf("xoje binary directory (%s) is correctly configured in PATH.", binDir)
	return res
}

func canonicalPath(path string) string {
	abs, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		abs = filepath.Clean(path)
	}
	if resolved, err := filepath.EvalSymlinks(abs); err == nil {
		return resolved
	}
	return abs
}

// NetworkCheck verifies HTTP connectivity to a target URL.
type NetworkCheck struct {
	Target string
}

func (c *NetworkCheck) ID() string   { return "network" }
func (c *NetworkCheck) Name() string { return "Network Connectivity" }
func (c *NetworkCheck) Run() Result {
	res := Result{Name: c.Name(), Status: StatusPass}
	client := &http.Client{Timeout: 3 * time.Second}
	response, err := client.Get(c.Target)
	if err != nil {
		res.Status = StatusFail
		res.Detail = fmt.Sprintf("Could not reach %s", c.Target)
		res.Remedy = "Check your internet connection or proxy settings."
		return res
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		res.Status = StatusFail
		res.Detail = fmt.Sprintf("%s returned HTTP %d", c.Target, response.StatusCode)
		res.Remedy = "Check your internet connection or proxy settings."
		return res
	}

	res.Detail = fmt.Sprintf("Successfully reached %s (HTTP 200)", c.Target)
	return res
}

// ConfigCheck verifies that the configuration file is readable and complete.
type ConfigCheck struct {
	Path string
}

func (c *ConfigCheck) ID() string   { return "config-file" }
func (c *ConfigCheck) Name() string { return "Configuration File" }
func (c *ConfigCheck) Run() Result {
	res := Result{Name: c.Name(), Status: StatusPass}
	data, err := os.ReadFile(c.Path)
	if err != nil {
		res.Status = StatusFail
		res.Detail = fmt.Sprintf("Could not read %s: %v", c.Path, err)
		res.Remedy = "Run 'xoje init' to create the configuration or fix its permissions."
		return res
	}

	var config struct {
		ActivePersona  string    `json:"active_persona"`
		InstalledTools *[]string `json:"installed_tools"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		res.Status = StatusFail
		res.Detail = fmt.Sprintf("Invalid JSON in %s: %v", c.Path, err)
		res.Remedy = "Repair the JSON configuration or recreate it with 'xoje init'."
		return res
	}
	if strings.TrimSpace(config.ActivePersona) == "" || config.InstalledTools == nil {
		res.Status = StatusFail
		res.Detail = "Configuration is missing required fields: active_persona and installed_tools."
		res.Remedy = "Add the required fields or recreate the configuration with 'xoje init'."
		return res
	}

	res.Detail = c.Path + " is valid."
	return res
}

// ToolsCheck verifies that all registered tools are still in PATH and functional.
type ToolsCheck struct {
	Tools []string
}

func (c *ToolsCheck) ID() string   { return "tools-registry" }
func (c *ToolsCheck) Name() string { return "Managed Tools" }
func (c *ToolsCheck) Run() Result {
	res := Result{Name: c.Name(), Status: StatusPass}
	if len(c.Tools) == 0 {
		res.Detail = "No tools registered in configuration."
		return res
	}

	var missing []string
	var collisions []string
	var okCount int

	// Map of tool names to their expected binary names
	binaryMap := make(map[string]string)
	for _, t := range install.Registry {
		binaryMap[t.Name] = t.BinaryName
	}

	for _, toolName := range c.Tools {
		binaryName, exists := binaryMap[toolName]
		if !exists {
			// Fallback to name if not in registry
			binaryName = toolName
		}

		// Use 'which -a' to detect multiple occurrences
		out, err := exec.Command("which", "-a", binaryName).Output()
		if err != nil {
			missing = append(missing, toolName)
			continue
		}

		paths := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(paths) > 1 {
			collisions = append(collisions, fmt.Sprintf("%s (%d paths)", toolName, len(paths)))
		}
		okCount++
	}

	if len(missing) > 0 {
		res.Status = StatusWarn
		res.Detail = fmt.Sprintf("%d/%d tools found. Missing: %s", okCount, len(c.Tools), strings.Join(missing, ", "))
		res.Remedy = "Run 'xoje update' to restore missing tools."
		res.RemedyAction = "update-all"
		return res
	}

	if len(collisions) > 0 {
		res.Status = StatusWarn
		res.Detail = fmt.Sprintf("All tools found, but COLLISIONS detected: %s", strings.Join(collisions, ", "))
		res.Remedy = "Manually remove redundant binaries from your PATH to avoid unexpected behavior."
		return res
	}

	res.Detail = fmt.Sprintf("All %d registered tools are available in PATH with no collisions.", len(c.Tools))
	return res
}

// EngramDoctorCheck verifies the internal health of Engram (sync status, locks, etc.)
type EngramDoctorCheck struct{}

func (c *EngramDoctorCheck) ID() string   { return "engram-doctor" }
func (c *EngramDoctorCheck) Name() string { return "Engram Memory Health" }
func (c *EngramDoctorCheck) Run() Result {
	res := Result{Name: c.Name(), Status: StatusPass}

	path, err := exec.LookPath("engram")
	if err != nil {
		res.Status = StatusFail
		res.Detail = "engram binary not found in PATH."
		res.Remedy = "Run 'xoje install engram' (part of gentle-ai ecosystem)."
		return res
	}

	// Run engram doctor --json to get structured internal health
	cmd := exec.Command(path, "doctor", "--json")
	output, _ := cmd.CombinedOutput()

	// If it fails or output is empty, we report unreachable as a Warn
	if len(output) == 0 {
		res.Status = StatusWarn
		res.Detail = "Engram internal doctor is unresponsive."
		res.Remedy = "Ensure Engram is installed and functional."
		return res
	}

	// Very basic check for "blocked" status in the JSON string
	// for a more robust approach we'd use json.Unmarshal
	outStr := string(output)
	if strings.Contains(outStr, "\"status\": \"blocked\"") || strings.Contains(outStr, "\"status\":\"blocked\"") {
		res.Status = StatusWarn
		res.Detail = "Engram sync is BLOCKED by pending mutations with missing fields."
		res.Remedy = "Run 'engram cloud upgrade doctor' and follow repair steps."
		return res
	}

	if strings.Contains(outStr, "\"status\": \"error\"") || strings.Contains(outStr, "\"status\":\"error\"") {
		res.Status = StatusFail
		res.Detail = "Engram internal state has ERRORS."
		res.Remedy = "Run 'engram doctor' manually to diagnose."
		return res
	}

	res.Detail = "Engram internal state and sync are healthy."
	return res
}

// GentleAIDoctorCheck triggers the native gentle-ai doctor and reports its status.
type GentleAIDoctorCheck struct{}

func (c *GentleAIDoctorCheck) ID() string   { return "gentle-ai-doctor" }
func (c *GentleAIDoctorCheck) Name() string { return "Gentle-AI Ecosystem Health" }
func (c *GentleAIDoctorCheck) Run() Result {
	res := Result{Name: c.Name(), Status: StatusPass}

	path, err := exec.LookPath("gentle-ai")
	if err != nil {
		res.Status = StatusFail
		res.Detail = "gentle-ai binary not found in PATH."
		res.Remedy = "Run 'xoje install gentle-ai' to restore the core ecosystem."
		res.RemedyAction = "install gentle-ai"
		return res
	}

	cmd := exec.Command(path, "doctor")
	output, err := cmd.CombinedOutput()

	// gentle-ai doctor returns non-zero exit code if unhealthy
	outStr := string(output)
	res.Detail = strings.TrimSpace(outStr)

	if err != nil {
		res.Status = StatusWarn // We mark as Warn so it doesn't block but signals issues
		if strings.Contains(outStr, "unhealthy") {
			res.Remedy = "Run 'gentle-ai doctor' for details and apply remedies."
		} else {
			res.Detail = fmt.Sprintf("Failed to run doctor: %v\n%s", err, outStr)
		}
	}

	return res
}

// DefaultChecks returns the standard suite of diagnostics for xoje.
func DefaultChecks(tools []string, configPath string) []Check {
	return []Check{
		&GoVersionCheck{Required: "1.22"},
		&ExecutableCheck{Binary: "node", Label: "Node", Remedy: "Install Node.js and ensure node is in PATH."},
		&ExecutableCheck{Binary: "git", Label: "Git", Remedy: "Install Git and ensure git is in PATH."},
		&EnvPathCheck{},
		&NetworkCheck{Target: "https://proxy.golang.org"},
		&NetworkCheck{Target: "https://github.com"},
		&ConfigCheck{Path: configPath},
		&ToolsCheck{Tools: tools},
		&EngramDoctorCheck{},
		&GentleAIDoctorCheck{},
	}
}
