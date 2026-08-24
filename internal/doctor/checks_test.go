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
