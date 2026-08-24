package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Xoje-Tech/xoje-environment/internal/config"
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
		_ = os.Remove(markerFile)
		err := Install("vim", &config.Config{})
		if err == nil || !strings.Contains(err.Error(), "tool not found") {
			t.Fatalf("expected whitelist rejection, got %v", err)
		}
		if _, statErr := os.Stat(markerFile); !os.IsNotExist(statErr) {
			t.Fatalf("installer command ran before whitelist rejection: stat err=%v", statErr)
		}
	})

	t.Run("Update all only visits registered tools and propagates failures", func(t *testing.T) {
		_ = os.Remove(markerFile)
		cfg := &config.Config{InstalledTools: []string{"not-in-whitelist"}}
		err := UpdateAll(cfg)
		if err == nil || !strings.Contains(err.Error(), "not-in-whitelist") {
			t.Fatalf("expected registered-tool failure, got %v", err)
		}
		if _, statErr := os.Stat(markerFile); !os.IsNotExist(statErr) {
			t.Fatalf("installer command ran for rejected registered tool: stat err=%v", statErr)
		}
	})

	t.Run("Normal install - arguments and GOBIN", func(t *testing.T) {
		os.Remove(markerFile)
		customBin := filepath.Join(tmpDir, "custom-bin")

		// Create a stub for gentle-ai in customBin to satisfy Post-install verification
		if err := os.MkdirAll(customBin, 0755); err != nil {
			t.Fatal(err)
		}
		gaStubPath := filepath.Join(customBin, "gentle-ai")
		gaMarker := filepath.Join(tmpDir, "ga_called.txt")
		gaStubContent := "#!/bin/sh\necho \"$@\" >> " + gaMarker + "\nexit 0"
		os.WriteFile(gaStubPath, []byte(gaStubContent), 0755)

		// Add customBin to PATH so IsInstalled(gentle-ai) finds it
		t.Setenv("PATH", customBin+string(os.PathListSeparator)+tmpDir+string(os.PathListSeparator)+oldPath)

		// Use gentle-ai from Registry
		var gaTool Tool
		for _, tool := range Registry {
			if tool.Name == "gentle-ai" {
				gaTool = tool
				break
			}
		}

		// Use isUpdate=false to force a fresh install in the stub env
		err := InstallTool(gaTool, customBin, false)
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

		gaTool := Registry[0] // gentle-ai
		err := InstallTool(gaTool, "", true)
		if err == nil || !strings.Contains(err.Error(), "error message") {
			t.Errorf("expected error containing 'error message', got %v", err)
		}
	})
}

func TestInstallRegistersAndPersistsTool(t *testing.T) {
	tmpDir := t.TempDir()
	goStub := filepath.Join(tmpDir, "go")
	if err := os.WriteFile(goStub, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", tmpDir)
	t.Setenv("HOME", tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")
	cfg, err := config.LoadOrCreate(configPath)
	if err != nil {
		t.Fatal(err)
	}

	if err := Install("gentle-ai", cfg); err != nil {
		t.Fatalf("Install() error = %v", err)
	}
	if len(cfg.InstalledTools) != 1 || cfg.InstalledTools[0] != "gentle-ai" {
		t.Fatalf("installed tools = %v, want [gentle-ai]", cfg.InstalledTools)
	}

	reloaded, err := config.LoadOrCreate(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(reloaded.InstalledTools) != 1 || reloaded.InstalledTools[0] != "gentle-ai" {
		t.Fatalf("persisted installed tools = %v, want [gentle-ai]", reloaded.InstalledTools)
	}
}

func TestUpdateRegistersAndPersistsTool(t *testing.T) {
	tmpDir := t.TempDir()
	goStub := filepath.Join(tmpDir, "go")
	if err := os.WriteFile(goStub, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", tmpDir)
	t.Setenv("HOME", tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")
	cfg, err := config.LoadOrCreate(configPath)
	if err != nil {
		t.Fatal(err)
	}

	if err := Update("gentle-ai", cfg); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	reloaded, err := config.LoadOrCreate(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(reloaded.InstalledTools) != 1 || reloaded.InstalledTools[0] != "gentle-ai" {
		t.Fatalf("persisted installed tools = %v, want [gentle-ai]", reloaded.InstalledTools)
	}
}
