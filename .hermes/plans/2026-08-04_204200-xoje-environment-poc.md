# xoje-environment Implementation Plan

> **For Hermes:** Use subagent-driven-development skill to implement this plan task-by-task.

**Goal:** Build a robust, platform-agnostic CLI + TUI tool (`xoje-environment`) in Go to manage, install, and sync our development environment (skills, identity, configurations) across different agent hosts, starting with a Proof of Concept that installs `gentle-ai`.

**Architecture:** A compiled Go binary with a clean, domain-driven structure. It splits responsibilities into distinct modules: `cli` (argument parsing, subcommands), `tui` (Bubbletea screens for interactive installation and status), `agents` (adaptor layer for projecting configurations), and `env` (environment detection and validation).

**Tech Stack:** Go (1.22+), Charm Bubbletea (TUI framework), YAML/JSON.

---

## Task Breakdown

### Task 1: Setup Workspace & Go Modules Scaffolding

**Objective:** Initialize the project workspace, set up Go modules, and create the core folder structure.

**Files:**
- Create: `/home/hermes/projects/xoje-environment/go.mod`
- Create: `/home/hermes/projects/xoje-environment/cmd/xoje/main.go`
- Create: `/home/hermes/projects/xoje-environment/.gitignore`

**Step 1: Write a minimal main.go**

Create `/home/hermes/projects/xoje-environment/cmd/xoje/main.go`:
```go
package main

import "fmt"

func main() {
	fmt.Println("xoje-environment v1.0.0 PoC")
}
```

Create `/home/hermes/projects/xoje-environment/go.mod`:
```go
module github.com/Xoje-Tech/xoje-environment

go 1.22
```

Create `/home/hermes/projects/xoje-environment/.gitignore`:
```gitignore
bin/
dist/
xoje
.env
```

**Step 2: Run build to verify setup**

Run: `go run /home/hermes/projects/xoje-environment/cmd/xoje/main.go`
Expected: `xoje-environment v1.0.0 PoC`

**Step 3: Commit**

```bash
git init /home/hermes/projects/xoje-environment
git add /home/hermes/projects/xoje-environment/
git commit -m "chore: initial repository structure and hello world"
```

---

### Task 2: Define Core Configuration & State (TDD)

**Objective:** Model how the environment tracks installed tools and active configurations inside a simple YAML/JSON file (`xoje-config.yaml`).

**Files:**
- Create: `/home/hermes/projects/xoje-environment/internal/config/config.go`
- Create: `/home/hermes/projects/xoje-environment/internal/config/config_test.go`

**Step 1: Write failing configuration test**

Create `/home/hermes/projects/xoje-environment/internal/config/config_test.go`:
```go
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func testTempConfigPath(t *testing.T) string {
	t.Helper()
	tmp, err := os.MkdirTemp("", "xoje-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	return filepath.Join(tmp, "xoje-config.yaml")
}

func TestLoadOrCreateConfig(t *testing.T) {
	path := testTempConfigPath(t)
	defer os.RemoveAll(filepath.Dir(path))

	// Test default creation
	cfg, err := LoadOrCreate(path)
	if err != nil {
		t.Fatalf("LoadOrCreate failed: %v", err)
	}

	if cfg.ActivePersona != "gandalf" {
		t.Errorf("expected default persona 'gandalf', got '%s'", cfg.ActivePersona)
	}

	if len(cfg.InstalledTools) != 0 {
		t.Errorf("expected 0 default tools, got %v", cfg.InstalledTools)
	}

	// Test update and save
	cfg.InstalledTools = append(cfg.InstalledTools, "gentle-ai")
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Reload and verify
	cfg2, err := LoadOrCreate(path)
	if err != nil {
		t.Fatalf("Reload failed: %v", err)
	}

	if len(cfg2.InstalledTools) != 1 || cfg2.InstalledTools[0] != "gentle-ai" {
		t.Errorf("reloaded config does not match saved config: %v", cfg2.InstalledTools)
	}
}
```

**Step 2: Run test to verify failure**

Run: `go test /home/hermes/projects/xoje-environment/internal/config/... -v`
Expected: FAIL — Compilation error (LoadOrCreate not defined)

**Step 3: Write minimal implementation**

Create `/home/hermes/projects/xoje-environment/internal/config/config.go` using a JSON encoder/decoder for zero external dependency compilation (or standard library encoders):
```go
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	ActivePersona  string   `json:"active_persona"`
	InstalledTools []string `json:"installed_tools"`
	path           string
}

func LoadOrCreate(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, err
		}
		cfg := &Config{
			ActivePersona:  "gandalf",
			InstalledTools: []string{},
			path:           path,
		}
		if err := cfg.Save(); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	cfg.path = path
	return &cfg, nil
}

func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0644)
}
```

**Step 4: Run test to verify pass**

Run: `go test /home/hermes/projects/xoje-environment/internal/config/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add /home/hermes/projects/xoje-environment/internal/config/
git commit -m "feat(config): implement JSON-backed config loader and saver"
```

---

### Task 3: Environment Diagnosis / Verification Engine (TDD)

**Objective:** Implement a diagnostics check layer to verify the host machine has prerequisite dependencies installed (`go`, `node`, `git`).

**Files:**
- Create: `/home/hermes/projects/xoje-environment/internal/verify/verify.go`
- Create: `/home/hermes/projects/xoje-environment/internal/verify/verify_test.go`

**Step 1: Write failing verification test**

Create `/home/hermes/projects/xoje-environment/internal/verify/verify_test.go`:
```go
package verify

import "testing"

func TestVerifyCommand(t *testing.T) {
	// Test a command we know exists
	res := VerifyExecutable("go")
	if !res.Available {
		t.Errorf("expected 'go' to be available, got false")
	}

	// Test a command we know doesn't exist
	res2 := VerifyExecutable("nonexistent-command-xyz")
	if res2.Available {
		t.Errorf("expected 'nonexistent-command-xyz' to be unavailable, got true")
	}
}
```

**Step 2: Run test to verify failure**

Run: `go test /home/hermes/projects/xoje-environment/internal/verify/... -v`
Expected: FAIL — Compilation error (VerifyExecutable not defined)

**Step 3: Write minimal implementation**

Create `/home/hermes/projects/xoje-environment/internal/verify/verify.go`:
```go
package verify

import "os/exec"

type DiagnosticResult struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
	Path      string `json:"path"`
}

func VerifyExecutable(name string) DiagnosticResult {
	path, err := exec.LookPath(name)
	if err != nil {
		return DiagnosticResult{
			Name:      name,
			Available: false,
			Path:      "",
		}
	}
	return DiagnosticResult{
		Name:      name,
		Available: true,
		Path:      path,
	}
}
```

**Step 4: Run test to verify pass**

Run: `go test /home/hermes/projects/xoje-environment/internal/verify/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add /home/hermes/projects/xoje-environment/internal/verify/
git commit -m "feat(verify): implement LookPath-based diagnostic check"
```

---

### Task 4: Implement Installer Engine (TDD)

**Objective:** Build the code that runs the `go install` command to install `gentle-ai` on the local machine and updates the tool registry.

**Files:**
- Create: `/home/hermes/projects/xoje-environment/internal/install/install.go`
- Create: `/home/hermes/projects/xoje-environment/internal/install/install_test.go`

**Step 1: Write failing installer test**

Create `/home/hermes/projects/xoje-environment/internal/install/install_test.go`:
```go
package install

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallGentleAI(t *testing.T) {
	// Mock configuration
	tmp, err := os.MkdirTemp("", "xoje-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmp)

	testBinDir := filepath.Join(tmp, "bin")
	err = os.MkdirAll(testBinDir, 0755)
	if err != nil {
		t.Fatalf("failed to create temp bin: %v", err)
	}

	// Run installation using a dry run or mocked install function to verify logic
	err = InstallTool("gentle-ai", testBinDir, true) // dry run
	if err != nil {
		t.Fatalf("InstallTool failed: %v", err)
	}
}
```

**Step 2: Run test to verify failure**

Run: `go test /home/hermes/projects/xoje-environment/internal/install/... -v`
Expected: FAIL — Compilation error (InstallTool not defined)

**Step 3: Write minimal implementation**

Create `/home/hermes/projects/xoje-environment/internal/install/install.go`:
```go
package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func InstallTool(name string, customBinDir string, dryRun bool) error {
	if name != "gentle-ai" {
		return fmt.Errorf("unknown tool: %s", name)
	}

	if dryRun {
		return nil
	}

	// Execute: go install github.com/gentleman-programming/gentle-ai/v2/cmd/gentle-ai@latest
	cmd := exec.Command("go", "install", "github.com/gentleman-programming/gentle-ai/v2/cmd/gentle-ai@latest")
	
	// Override GOBIN to point to the desired location if supplied
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
```

**Step 4: Run test to verify pass**

Run: `go test /home/hermes/projects/xoje-environment/internal/install/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add /home/hermes/projects/xoje-environment/internal/install/
git commit -m "feat(install): add engine to execute tool installation via go install"
```

---

### Task 5: Implement CLI Arguments Parsing (TDD)

**Objective:** Standardize argument parsing to support commands: `xoje diagnose`, `xoje install gentle-ai`, and `xoje tui` (default).

**Files:**
- Create: `/home/hermes/projects/xoje-environment/internal/cli/cli.go`
- Create: `/home/hermes/projects/xoje-environment/internal/cli/cli_test.go`

**Step 1: Write failing CLI routing test**

Create `/home/hermes/projects/xoje-environment/internal/cli/cli_test.go`:
```go
package cli

import "testing"

func TestParseArgs(t *testing.T) {
	cmd, tool, err := Parse([]string{"diagnose"})
	if err != nil || cmd != "diagnose" {
		t.Errorf("expected 'diagnose', got %s (%v)", cmd, err)
	}

	cmd2, tool2, err := Parse([]string{"install", "gentle-ai"})
	if err != nil || cmd2 != "install" || tool2 != "gentle-ai" {
		t.Errorf("expected install gentle-ai, got %s %s (%v)", cmd2, tool2, err)
	}

	cmd3, _, err := Parse([]string{})
	if err != nil || cmd3 != "tui" {
		t.Errorf("expected default tui, got %s (%v)", cmd3, err)
	}
}
```

**Step 2: Run test to verify failure**

Run: `go test /home/hermes/projects/xoje-environment/internal/cli/... -v`
Expected: FAIL — Compilation error (Parse not defined)

**Step 3: Write minimal implementation**

Create `/home/hermes/projects/xoje-environment/internal/cli/cli.go`:
```go
package cli

import "fmt"

func Parse(args []string) (string, string, error) {
	if len(args) == 0 {
		return "tui", "", nil
	}

	subcommand := args[0]
	switch subcommand {
	case "diagnose", "bootstrap":
		return subcommand, "", nil
	case "install":
		if len(args) < 2 {
			return "", "", fmt.Errorf("usage: xoje install [tool_name]")
		}
		return "install", args[1], nil
	case "tui":
		return "tui", "", nil
	default:
		return "", "", fmt.Errorf("unknown subcommand: %s", subcommand)
	}
}
```

**Step 4: Run test to verify pass**

Run: `go test /home/hermes/projects/xoje-environment/internal/cli/... -v`
Expected: PASS

**Step 5: Commit**

```bash
git add /home/hermes/projects/xoje-environment/internal/cli/
git commit -m "feat(cli): add lightweight CLI subcommand parser"
```

---

### Task 6: Setup TUI Bubbletea Program Shell (TDD)

**Objective:** Scaffold the interactive TUI shell using Bubbletea, adding basic keyboard navigation (Up/Down, Enter, Esc/q to quit).

**Files:**
- Create: `/home/hermes/projects/xoje-environment/internal/tui/tui.go`
- Create: `/home/hermes/projects/xoje-environment/internal/tui/tui_test.go`

**Step 1: Add Bubbletea dependency to go.mod**

Add `github.com/charmbracelet/bubbletea` to the module configuration or via download.

**Step 2: Write failing TUI rendering test**

Create `/home/hermes/projects/xoje-environment/internal/tui/tui_test.go`:
```go
package tui

import (
	"strings"
	"testing"
)

func TestInitialModelView(t *testing.T) {
	m := NewModel()
	view := m.View()

	if !strings.Contains(view, "xoje-environment") {
		t.Errorf("expected view to contain header 'xoje-environment', got:\n%s", view)
	}
}
```

**Step 3: Run test to verify failure**

Run: `go test /home/hermes/projects/xoje-environment/internal/tui/... -v`
Expected: FAIL — Compilation error (NewModel not defined)

**Step 4: Fetch Charm libraries & write minimal Bubbletea implementation**

Run:
```bash
go get github.com/charmbracelet/bubbletea@v1.0.0
```

Create `/home/hermes/projects/xoje-environment/internal/tui/tui.go`:
```go
package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	cursor int
	choices []string
	selected map[int]struct{}
}

func NewModel() Model {
	return Model{
		choices:  []string{"diagnose", "install gentle-ai", "exit"},
		selected: make(map[int]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			if m.cursor == 2 {
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	var s strings.Builder
	s.WriteString("=================================\n")
	s.WriteString("  xoje-environment CLI & TUI  \n")
	s.WriteString("=================================\n\n")

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}

	s.WriteString("\nPress q/ctrl+c to exit.\n")
	return s.String()
}
```

**Step 5: Run test to verify pass**

Run: `go test /home/hermes/projects/xoje-environment/internal/tui/... -v`
Expected: PASS

**Step 6: Commit**

```bash
git add /home/hermes/projects/xoje-environment/go.mod /home/hermes/projects/xoje-environment/go.sum /home/hermes/projects/xoje-environment/internal/tui/
git commit -m "feat(tui): setup Bubbletea state machine and primary rendering engine"
```

---

### Task 7: Wire everything together in Main entrypoint

**Objective:** Read configuration, parse arguments, and run the diagnostics engine or the TUI based on the subcommand passed.

**Files:**
- Modify: `/home/hermes/projects/xoje-environment/cmd/xoje/main.go`

**Step 1: Write integration logic in main.go**

Modify `/home/hermes/projects/xoje-environment/cmd/xoje/main.go` to connect CLI parser, diagnostic executor, installer, and Bubbletea TUI:
```go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Xoje-Tech/xoje-environment/internal/cli"
	"github.com/Xoje-Tech/xoje-environment/internal/config"
	"github.com/Xoje-Tech/xoje-environment/internal/install"
	"github.com/Xoje-Tech/xoje-environment/internal/tui"
	"github.com/Xoje-Tech/xoje-environment/internal/verify"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error resolving home directory: %v\n", err)
		os.Exit(1)
	}

	configPath := filepath.Join(home, ".config", "xoje", "config.json")
	cfg, err := config.LoadOrCreate(configPath)
	if err != nil {
		fmt.Printf("Error initializing configuration: %v\n", err)
		os.Exit(1)
	}

	cmd, tool, err := cli.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Arguments error: %v\n", err)
		os.Exit(1)
	}

	switch cmd {
	case "diagnose", "bootstrap":
		fmt.Println("Running xoje diagnostics...")
		checks := []string{"go", "node", "git"}
		allPassed := true
		for _, name := range checks {
			res := verify.VerifyExecutable(name)
			if res.Available {
				fmt.Printf("✅ [%s] is available at %s\n", name, res.Path)
			} else {
				fmt.Printf("❌ [%s] is NOT found on the system path\n", name)
				allPassed = false
			}
		}
		if !allPassed {
			fmt.Println("Warning: some prerequisites are missing. Go installation may fail.")
		}
	case "install":
		fmt.Printf("Installing tool: %s...\n", tool)
		err := install.InstallTool(tool, "", false)
		if err != nil {
			fmt.Printf("❌ Installation failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ %s successfully installed!\n", tool)

		// Record in active config
		cfg.InstalledTools = append(cfg.InstalledTools, tool)
		if err := cfg.Save(); err != nil {
			fmt.Printf("Failed to update config: %v\n", err)
		}
	case "tui":
		p := tea.NewProgram(tui.NewModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("TUI Error: %v\n", err)
			os.Exit(1)
		}
	}
}
```

**Step 2: Run Go test to ensure everything compiles**

Run: `go test /home/hermes/projects/xoje-environment/...`
Expected: PASS (all individual test suites passed)

**Step 3: Commit**

```bash
git add /home/hermes/projects/xoje-environment/cmd/xoje/main.go
git commit -m "feat(main): integrate config, cli, verify and tui into main executable"
```

---

## E2E Validation Phase

Once everything compiles cleanly, verify the PoC functionality end-to-end:

1. **Build the binary:**
   ```bash
   cd /home/hermes/projects/xoje-environment
   go build -o bin/xoje cmd/xoje/main.go
   ```
2. **Execute diagnosis subcommand:**
   ```bash
   ./bin/xoje diagnose
   ```
   *Verify it output checks for `go`, `node`, `git`.*
3. **Execute tool installation:**
   ```bash
   ./bin/xoje install gentle-ai
   ```
   *Verify that the CLI installs gentle-ai successfully via `go install`.*
4. **Execute `gentle-ai` to ensure the installation succeeded:**
   ```bash
   gentle-ai version
   ```
   *Expected: displays gentle-ai's compiled version.*
