package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var whitelist = map[string]string{
	"gentle-ai": "github.com/gentleman-programming/gentle-ai/v2/cmd/gentle-ai@latest",
}

func InstallTool(name string, customBinDir string, dryRun bool) error {
	module, ok := whitelist[name]
	if !ok {
		return fmt.Errorf("unknown tool: %s", name)
	}

	if dryRun {
		return nil
	}

	cmd := exec.Command("go", "install", module)

	if customBinDir != "" {
		absPath, err := filepath.Abs(customBinDir)
		if err == nil {
			cmd.Env = append(os.Environ(), "GOBIN="+absPath)
		}
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go install failed: %w (output: %s)", err, string(output))
	}

	return nil
}
