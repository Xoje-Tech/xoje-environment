package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallTool(t *testing.T) {
	tmpDir := t.TempDir()
	stubPath := filepath.Join(tmpDir, "go")
	markerFile := filepath.Join(tmpDir, "called.txt")

	// Create a stub 'go' binary (shell script)
	stubContent := "#!/bin/sh\necho \"$@\" > " + markerFile + "\necho \"GOBIN=$GOBIN\" >> " + markerFile + "\nexit 0"
	if err := os.WriteFile(stubPath, []byte(stubContent), 0755); err != nil {
		t.Fatalf("failed to create stub: %v", err)
	}

	// Setup PATH to include our stub
	oldPath := os.Getenv("PATH")
	t.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath)
	// Clear GOBIN for tests
	t.Setenv("GOBIN", "")

	t.Run("Whitelist check - unknown tool", func(t *testing.T) {
		err := InstallTool("vim", "", false)
		if err == nil || !strings.Contains(err.Error(), "unknown tool") {
			t.Errorf("expected unknown tool error, got %v", err)
		}
	})

	t.Run("Dry run - no execution", func(t *testing.T) {
		os.Remove(markerFile)
		err := InstallTool("gentle-ai", "", true)
		if err != nil {
			t.Errorf("expected no error in dry run, got %v", err)
		}
		if _, err := os.Stat(markerFile); !os.IsNotExist(err) {
			t.Error("stub was called in dry run, but it should not have been")
		}
	})

	t.Run("Normal install - arguments and GOBIN", func(t *testing.T) {
		os.Remove(markerFile)
		customBin := filepath.Join(tmpDir, "custom-bin")
		err := InstallTool("gentle-ai", customBin, false)
		if err != nil {
			t.Fatalf("InstallTool failed: %v", err)
		}

		data, _ := os.ReadFile(markerFile)
		output := string(data)

		if !strings.Contains(output, "install github.com/gentleman-programming/gentle-ai/v2/cmd/gentle-ai@latest") {
			t.Errorf("stub called with wrong args: %s", output)
		}

		absBin, _ := filepath.Abs(customBin)
		if !strings.Contains(output, "GOBIN="+absBin) {
			t.Errorf("GOBIN not set correctly in stub env: %s", output)
		}
	})

	t.Run("Failure reporting", func(t *testing.T) {
		// Make stub fail
		stubContentFail := "#!/bin/sh\necho \"error message\"\nexit 1"
		os.WriteFile(stubPath, []byte(stubContentFail), 0755)

		err := InstallTool("gentle-ai", "", false)
		if err == nil || !strings.Contains(err.Error(), "error message") {
			t.Errorf("expected error containing 'error message', got %v", err)
		}
	})
}
