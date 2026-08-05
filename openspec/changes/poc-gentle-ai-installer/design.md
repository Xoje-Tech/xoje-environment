# Design: xoje-environment PoC — gentle-ai installer

## Technical Approach

Greenfield Go 1.22+ CLI+TUI (module `github.com/Xoje-Tech/xoje-environment`). Five decoupled internal packages (`config`, `verify`, `install`, `cli`, `tui`) plus `cmd/xoje` wiring, per the approved plan (Tasks 1–7) and the five capability specs. Only external dependency: Bubbletea v1.0.0; stdlib JSON for config. Strict TDD (`go test ./...`). Entry flow: `config.LoadOrCreate` → `cli.Parse` → dispatch to diagnose/install/tui.

## Architecture Decisions

| # | Decision | Alternatives | Rationale |
|---|----------|--------------|-----------|
| D1 | Domain packages `internal/{config,verify,install,cli,tui}` + `cmd/xoje` | Flat `main`; `agents/`+`env/` layering | Plan/proposal mandate; per-package TDD; one-way deps (main→cli→capabilities); `agents/` out of PoC scope |
| D2 | stdlib JSON config (`active_persona`, `installed_tools`) | YAML (gopkg.in/yaml.v3) | config-state REQUIRES JSON snake_case; zero extra deps; plan uses `json.MarshalIndent`; temp-path tests kill format-drift risk |
| D3 | `exec.LookPath` verification + aggregate readiness | Manual PATH parsing; version probes | environment-diagnosis mandates LookPath; no subprocess in diagnosis; missing tools derivable from result list |
| D4 | Installer: whitelist gate → dry-run short-circuit → GOBIN override → CombinedOutput | Runner abstraction; unconditional exec | tool-install R1–R3; whitelist rejects unknowns before any exec; GOBIN `filepath.Abs`-normalized; output captured into errors |
| D5 | Registry update in `main` wiring, not `install` | `install` imports `config` | Keeps packages decoupled (config.yaml apply rules); matches plan Task 7 |
| D6 | Bubbletea `Model`/`Update`/`View`; cursor clamp; q/esc/ctrl+c quit | lipgloss styling; list library | interactive-tui: header, cursor, clamp, quit keys, Enter-on-exit; minimal, unit-testable |
| D7 | `cli.Parse(args) (command, tool, error)`; default `tui`; `bootstrap` alias; exit 0/1 in main | cobra | cli-commands; three subcommands don't justify cobra; alias + usage error + exit codes per spec |

## Data Flow

```
os.Args ──► cli.Parse ──► dispatch (main)
              │
config.LoadOrCreate(~/.config/xoje/config.json) ◄── first, always
              │
  diagnose ──► verify.VerifyExecutable(go|node|git) ──► report + readiness
  install  ──► install.InstallTool(tool, bin, dry) ──► main appends tool ──► config.Save
  tui      ──► tea.NewProgram(tui.NewModel()).Run()
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `go.mod`, `go.sum` | Create | module + `bubbletea@v1.0.0` |
| `.gitignore` | Create | `bin/`, `dist/`, `xoje`, `.env` |
| `cmd/xoje/main.go` | Create | wiring: config → parse → dispatch; exit codes |
| `internal/config/{config.go,config_test.go}` | Create | load-or-create, JSON save |
| `internal/verify/{verify.go,verify_test.go}` | Create | LookPath checks + aggregate |
| `internal/install/{install.go,install_test.go}` | Create | whitelist, dry-run, GOBIN, output capture |
| `internal/cli/{cli.go,cli_test.go}` | Create | subcommand routing |
| `internal/tui/{tui.go,tui_test.go}` | Create | Bubbletea model/update/view |

## Interfaces / Contracts

```go
// internal/config
type Config struct {
	ActivePersona  string   `json:"active_persona"`
	InstalledTools []string `json:"installed_tools"`
	path string // unexported
}
func LoadOrCreate(path string) (*Config, error)
func DefaultConfigPath() (string, error) // $HOME/.config/xoje/config.json
func (c *Config) Save() error

// internal/verify
type DiagnosticResult struct{ Name, Path string; Available bool }
func VerifyExecutable(name string) DiagnosticResult
func DiagnoseAll() ([]DiagnosticResult, bool) // results + all-available readiness

// internal/install
var whitelist = map[string]string{ // name → module@version
	"gentle-ai": "github.com/gentleman-programming/gentle-ai/v2/cmd/gentle-ai@latest",
}
func InstallTool(name, customBinDir string, dryRun bool) error

// internal/cli
func Parse(args []string) (command, tool string, err error)

// internal/tui
type Model struct{ cursor int; choices []string }
func NewModel() Model
func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m Model) View() string
```

## Sequence Diagram — install flow

```
User  main        config      cli        install      go (subprocess)
 │     │             │          │            │             │
 │ xoje install gentle-ai      │            │             │
 ├────►│             │          │            │             │
 │     ├─LoadOrCreate(config.json)           │             │
 │     │◄──── Config ──┤         │            │             │
 │     ├─Parse(["install","gentle-ai"])      │             │
 │     │◄── ("install","gentle-ai",nil)      │             │
 │     ├─InstallTool("gentle-ai","",false)──►│             │
 │     │             │          │  whitelist OK           │
 │     │             │          ├─go install …@latest────►│
 │     │             │          │◄──────── exit 0 ────────┤
 │     │◄──── nil ───┤          │             │            │
 │     ├─append "gentle-ai"; Save()           │            │
 │     │  (dry-run: returns after whitelist; no exec, no write)
```

## Testing Strategy

| Layer | What | How |
|-------|------|-----|
| Unit | config: create-defaults, load, round-trip, corrupt JSON, default path | `os.MkdirTemp` temp dirs; table tests |
| Unit | verify: present/absent LookPath, aggregate readiness | known name + `nonexistent-command-xyz` |
| Unit | install: whitelist accept/reject, dry-run no-op, GOBIN env, failure output capture | stub `go` on PATH via `t.Setenv` (no toolchain needed) |
| Unit | cli: default tui, diagnose/bootstrap, install ±tool, unknown | table tests on `Parse` |
| Unit | tui: header render, cursor move/clamp, quit keys, Enter-on-exit | `View`/`Update` assertions on model |
| E2E | build, diagnose, install, `gentle-ai version`, TUI q/ctrl+c | plan E2E phase; needs Go 1.22+ & network on host |

## Threat Matrix

| Boundary | Applicability | Design response | RED tests |
|---|---|---|---|
| Documentation-like paths | N/A — no doc-file execution; install target is a fixed whitelisted module path, never user-supplied | — | — |
| Git repository selection | N/A — tool never invokes git | — | — |
| Commit state | N/A — no commit automation | — | — |
| Push state | N/A — no push automation | — | — |
| PR commands | N/A — no PR automation | — | — |

The only applicable boundary is subprocess execution, governed by D4: whitelist gate rejects unknown tools before any exec; dry-run never execs; GOBIN is absolutized; subprocess output is captured into returned errors. RED tests map to tool-install scenarios (unknown tool → error + no exec; dry-run success without exec/write; GOBIN absolute; failure includes output).

## Migration / Rollout

No migration — greenfield. Rollback: delete `openspec/changes/poc-gentle-ai-installer/` + generated Go files; `git revert` if repo initialized (proposal).

## Open Questions

- None blocking. Note: E2E install/verify requires the Go 1.22+ toolchain + network on the host (not installed — apply phase must provision it first).
