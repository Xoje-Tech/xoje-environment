package doctor

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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

// EnvPathCheck verifies that Go binary directories are in the system PATH.
type EnvPathCheck struct{}

func (c *EnvPathCheck) ID() string   { return "env-path" }
func (c *EnvPathCheck) Name() string { return "System PATH Configuration" }
func (c *EnvPathCheck) Run() Result {
	res := Result{Name: c.Name(), Status: StatusPass}
	
	pathEnv := os.Getenv("PATH")
	paths := filepath.SplitList(pathEnv)
	
	goBin := os.Getenv("GOBIN")
	if goBin == "" {
		goPath := os.Getenv("GOPATH")
		if goPath == "" {
			home, _ := os.UserHomeDir()
			goPath = filepath.Join(home, "go")
		}
		goBin = filepath.Join(goPath, "bin")
	}

	absGoBin, _ := filepath.Abs(goBin)
	
	found := false
	for _, p := range paths {
		absP, _ := filepath.Abs(p)
		if absP == absGoBin {
			found = true
			break
		}
	}

	if !found {
		res.Status = StatusFail
		res.Detail = fmt.Sprintf("Go binary directory (%s) is NOT in your PATH.", goBin)
		res.Remedy = fmt.Sprintf("Add 'export PATH=$PATH:%s' to your shell profile.", goBin)
		return res
	}

	res.Detail = fmt.Sprintf("Go binary directory (%s) is correctly configured in PATH.", goBin)
	return res
}

// NetworkCheck verifies connectivity to a target host.
type NetworkCheck struct {
	Target string
}

func (c *NetworkCheck) ID() string   { return "network" }
func (c *NetworkCheck) Name() string { return "Network Connectivity" }
func (c *NetworkCheck) Run() Result {
	res := Result{Name: c.Name(), Status: StatusPass}
	conn, err := net.DialTimeout("tcp", c.Target+":80", 2*time.Second)
	if err != nil {
		res.Status = StatusFail
		res.Detail = "Could not reach " + c.Target
		res.Remedy = "Check your internet connection or proxy settings."
		return res
	}
	conn.Close()
	res.Detail = "Connected to " + c.Target
	return res
}

// ConfigCheck verifies the existence of a configuration file.
type ConfigCheck struct {
	Path string
}

func (c *ConfigCheck) ID() string   { return "config-file" }
func (c *ConfigCheck) Name() string { return "Configuration File" }
func (c *ConfigCheck) Run() Result {
	res := Result{Name: c.Name(), Status: StatusPass}
	_, err := os.Stat(c.Path)
	if os.IsNotExist(err) {
		res.Status = StatusFail
		res.Detail = "Missing " + c.Path
		res.Remedy = "Run 'xoje init' to generate a default configuration."
		return res
	}
	res.Detail = c.Path + " found."
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

	for _, tool := range c.Tools {
		// Use 'which -a' to detect multiple occurrences
		out, err := exec.Command("which", "-a", tool).Output()
		if err != nil {
			missing = append(missing, tool)
			continue
		}

		paths := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(paths) > 1 {
			collisions = append(collisions, fmt.Sprintf("%s (%d paths)", tool, len(paths)))
		}
		okCount++
	}

	if len(missing) > 0 {
		res.Status = StatusWarn
		res.Detail = fmt.Sprintf("%d/%d tools found. Missing: %s", okCount, len(c.Tools), strings.Join(missing, ", "))
		res.Remedy = "Run 'xoje update' to restore missing tools."
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
