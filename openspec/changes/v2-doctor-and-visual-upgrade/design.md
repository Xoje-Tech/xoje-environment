# Design: xoje-environment v2.0 — doctor & visual upgrade

## Technical Approach

The v2.0 upgrade refactors the core diagnostic logic into a modular, extensible `internal/doctor` package and introduces a unified styling layer in `internal/tui/styles`. We replace the primitive `diagnose` logic with an interface-driven engine that supports concurrent health checks and structured reporting with actionable remedies. The visual layer uses `charmbracelet/lipgloss` to implement a Rose Pine-inspired theme, applying it consistently across CLI outputs (using frames/borders) and TUI components.

## Architecture Decisions

### Decision: Interface-based Diagnostic Checks

**Choice**: Define a `Check` interface in `internal/doctor`.
**Alternatives considered**: Static helper functions (v1.0 approach), command-based registry.
**Rationale**: Interfaces allow for diverse check types (version, network, config) while keeping the Runner decoupled. It facilitates unit testing with mocks and future extensibility without modifying the orchestration logic.

### Decision: Centralized Styling Package

**Choice**: `internal/tui/styles` package containing Lipgloss style constants.
**Alternatives considered**: Inlining styles in components, global CSS-like configuration file.
**Rationale**: Centralizing styles ensures visual consistency between CLI commands and TUI views. It simplifies theme updates and reduces boilerplate in component code.

### Decision: Concurrent Doctor Execution

**Choice**: Execute checks in parallel using goroutines and `sync.WaitGroup` or channels.
**Alternatives considered**: Sequential execution.
**Rationale**: Network and file system probes can be slow. Parallel execution minimizes total wait time while maintaining a responsive CLI/TUI experience.

## Data Flow

The `doctor.Runner` collects registered `Check` implementations, executes them concurrently, and aggregates `Result` objects.

    main/tui ──→ doctor.Run() ──→ [Runner]
                                     │
               ┌──────────┬──────────┼──────────┐
               │          │          │          │
         VersionCheck NetworkCheck ConfigCheck Other...
               │          │          │          │
               └──────────┴──────────┼──────────┘
                                     │
    Result List ←────────────────────┘
         │
    styles.RenderReport() ──→ Output

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/tui/styles/styles.go` | Create | Lipgloss constants (Rose Pine) and framing functions. |
| `internal/doctor/doctor.go` | Create | `Check` interface, `Result` struct, and concurrent `Runner`. |
| `internal/doctor/checks.go` | Create | Concrete implementations: `GoVersionCheck`, `NetworkCheck`, `ConfigCheck`. |
| `internal/verify/verify.go` | Delete | Logic superseded by the `doctor` package. |
| `internal/tui/tui.go` | Modify | Apply Lipgloss styles to menu and update options (diagnose -> doctor). |
| `cmd/xoje/main.go` | Modify | Refactor `dispatch` to use `doctor.Run` and `styles` for output. |

## Interfaces / Contracts

### Doctor Check Interface
```go
package doctor

type Status int

const (
    StatusPass Status = iota
    StatusWarn
    StatusFail
)

type Result struct {
    Name     string
    Status   Status
    Detail   string
    Remedy   string
}

type Check interface {
    ID() string
    Name() string
    Run() Result
}
```

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | `doctor.Runner` | Mock `Check` implementations to verify aggregation and timeout logic. |
| Unit | Individual `Check` implementations | Subprocess/Network mocking (using `httptest`) to verify status logic. |
| Integration | `styles` rendering | Verify Lipgloss output strings contain expected color codes/borders (visual inspection + string match). |

## Threat Matrix

| Boundary | Minimum adversarial cases | Applicability | Design response | Planned RED tests |
|---|---|---|---|---|
| Documentation-like paths | `README.sh`, executable Markdown | N/A | No execution of user-controlled document paths in this scope. | N/A |
| Git repository selection | `git -C` | N/A | Doctor only checks git existence, not repo-specific routing. | N/A |
| Commit state | staged, empty index | N/A | No git state mutations. | N/A |
| Push state | tracking branch | N/A | No network pushes. | N/A |
| PR commands | explicit `--head` | N/A | No PR automation. | N/A |

*Note: While `doctor` probes network and runs `go version` (subprocess), these are fixed, non-user-controlled inputs in this phase.*

## Migration / Rollout

- The `diagnose` command will be aliased to `doctor` for one version before full removal.
- Existing `internal/verify` logic will be migrated to `internal/doctor` before the package is deleted.

## Open Questions

- [ ] Should we support JSON output for `xoje doctor` to allow machine readability? (Rationale: Good for CI, but keeping TUI/CLI focus for now).
- [ ] Do we need a timeout per-check or a global timeout for the Runner? (Decision: Per-check timeout is safer for slow networks).
