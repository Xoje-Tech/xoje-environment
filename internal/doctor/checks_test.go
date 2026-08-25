package doctor

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoVersionCheck(t *testing.T) {
	check := &GoVersionCheck{Required: "1.22"}
	result := check.Run()

	if result.Name != "Go Version" {
		t.Errorf("expected name 'Go Version', got %s", result.Name)
	}
}

func TestNetworkCheck(t *testing.T) {
	t.Run("200 response within timeout passes", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		result := (&NetworkCheck{Target: server.URL}).Run()
		if result.Name != "Network Connectivity" {
			t.Errorf("expected name 'Network Connectivity', got %s", result.Name)
		}
		if result.Status != StatusPass {
			t.Fatalf("status = %q, want %q: %#v", result.Status, StatusPass, result)
		}
	})

	t.Run("non-200 response fails with remedy", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer server.Close()

		result := (&NetworkCheck{Target: server.URL}).Run()
		if result.Status != StatusFail || result.Remedy == "" {
			t.Fatalf("expected failure with remedy: %#v", result)
		}
	})
}

func TestConfigCheck(t *testing.T) {
	t.Run("valid config passes", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "config.json")
		if err := os.WriteFile(path, []byte(`{"active_persona":"gandalf","installed_tools":[]}`), 0o600); err != nil {
			t.Fatal(err)
		}
		result := (&ConfigCheck{Path: path}).Run()
		if result.Name != "Configuration File" || result.Status != StatusPass {
			t.Fatalf("expected valid configuration to pass: %#v", result)
		}
	})

	for name, content := range map[string]string{
		"invalid JSON":           `{`,
		"missing active persona": `{"installed_tools":[]}`,
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.json")
			if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
				t.Fatal(err)
			}
			result := (&ConfigCheck{Path: path}).Run()
			if result.Status != StatusFail || result.Remedy == "" {
				t.Fatalf("expected invalid configuration to fail with remedy: %#v", result)
			}
		})
	}
}

func TestEnvPathCheckCanonicalBinDir(t *testing.T) {
	home := t.TempDir()
	canonicalBin := filepath.Join(home, ".local", "bin")
	if err := os.MkdirAll(canonicalBin, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	t.Setenv("GOBIN", filepath.Join(home, "ignored-gobin"))
	t.Setenv("GOPATH", filepath.Join(home, "ignored-gopath"))

	t.Run("missing canonical directory fails with actionable remedy", func(t *testing.T) {
		t.Setenv("PATH", "/usr/bin:/bin")
		result := (&EnvPathCheck{}).Run()
		if result.Status != StatusFail {
			t.Fatalf("status = %q, want %q", result.Status, StatusFail)
		}
		if !strings.Contains(result.Detail, canonicalBin) || !strings.Contains(result.Remedy, canonicalBin) {
			t.Fatalf("result should identify canonical bin dir %q: %#v", canonicalBin, result)
		}
	})

	t.Run("clean and trailing slash entries pass", func(t *testing.T) {
		t.Setenv("PATH", "/usr/bin"+string(os.PathListSeparator)+canonicalBin+string(os.PathSeparator))
		result := (&EnvPathCheck{}).Run()
		if result.Status != StatusPass {
			t.Fatalf("status = %q, want %q: %#v", result.Status, StatusPass, result)
		}
	})

	t.Run("symlink to canonical directory passes", func(t *testing.T) {
		link := filepath.Join(home, "bin-link")
		if err := os.Symlink(canonicalBin, link); err != nil {
			t.Fatal(err)
		}
		t.Setenv("PATH", link)
		result := (&EnvPathCheck{}).Run()
		if result.Status != StatusPass {
			t.Fatalf("status = %q, want %q: %#v", result.Status, StatusPass, result)
		}
	})
}

func TestDefaultChecksIncludeRuntimeReadiness(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	checks := DefaultChecks(nil, filepath.Join(t.TempDir(), "config.json"))

	for _, id := range []string{"node-runtime", "git-runtime"} {
		var found Check
		for _, check := range checks {
			if check.ID() == id {
				found = check
				break
			}
		}
		if found == nil {
			t.Errorf("DefaultChecks() missing %q readiness check", id)
			continue
		}
		result := found.Run()
		if result.Status != StatusFail || result.Remedy == "" {
			t.Errorf("%s result = %#v, want failure with remedy", id, result)
		}
	}
}

func TestDefaultChecksRunFullPrerequisiteSuite(t *testing.T) {
	home := t.TempDir()
	binDir := filepath.Join(home, ".local", "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, output := range map[string]string{
		"go":        "go version go1.22.0 linux/amd64",
		"node":      "v22.0.0",
		"git":       "git version 2.46.0",
		"engram":    `{"status":"ok"}`,
		"gentle-ai": "healthy",
	} {
		path := filepath.Join(binDir, name)
		content := "#!/bin/sh\nprintf '%s\\n' '" + output + "'\n"
		if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("HOME", home)
	t.Setenv("PATH", binDir)

	configPath := filepath.Join(home, ".config", "xoje", "config.json")
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, []byte(`{"active_persona":"gandalf","installed_tools":[]}`), 0o600); err != nil {
		t.Fatal(err)
	}

	checks := DefaultChecks(nil, configPath)
	wantIDs := []string{
		"go-version",
		"node-runtime",
		"git-runtime",
		"env-path",
		"network",
		"network",
		"config-file",
		"tools-registry",
		"engram-doctor",
		"gentle-ai-doctor",
	}
	wantNetworkTargets := []string{"https://proxy.golang.org", "https://github.com"}
	if len(checks) != len(wantIDs) {
		t.Fatalf("DefaultChecks() returned %d checks, want %d", len(checks), len(wantIDs))
	}
	for i, check := range checks {
		if check.ID() != wantIDs[i] {
			t.Errorf("DefaultChecks()[%d].ID() = %q, want %q", i, check.ID(), wantIDs[i])
		}
	}

	networkServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer networkServer.Close()
	for i, wantTarget := range wantNetworkTargets {
		networkCheck, ok := checks[4+i].(*NetworkCheck)
		if !ok {
			t.Fatalf("DefaultChecks()[%d] = %T, want *NetworkCheck", 4+i, checks[4+i])
		}
		if networkCheck.Target != wantTarget {
			t.Errorf("DefaultChecks()[%d].Target = %q, want %q", 4+i, networkCheck.Target, wantTarget)
		}
		networkCheck.Target = networkServer.URL
	}

	results := NewRunner(checks).Run()
	if len(results) != len(checks) {
		t.Fatalf("Runner.Run() returned %d results for %d checks", len(results), len(checks))
	}
	for i, result := range results {
		if result.Name != checks[i].Name() {
			t.Errorf("result[%d].Name = %q, want %q", i, result.Name, checks[i].Name())
		}
		if result.Detail == "" {
			t.Errorf("result[%d] (%s) has no reported detail", i, checks[i].ID())
		}
	}
	for _, index := range []int{4, 5} {
		if !strings.Contains(results[index].Detail, networkServer.URL) {
			t.Errorf("network result does not identify its exercised target: %#v", results[index])
		}
	}
}
