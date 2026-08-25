# Proposal: xoje-environment v2.0 — doctor & visual upgrade

## Intent

Upgrade `xoje-environment` from a PoC to a professional orchestration tool. Replace the basic `diagnose` command with a comprehensive `doctor` utility inspired by `gentle-ai`, and introduce a polished visual identity using the Lipgloss framework.

## Scope

### In Scope
- **`internal/doctor` Package**: Modular health check engine.
  - Interface-based checks (Version, Connectivity, Config, Disk).
  - Detailed results with statuses (PASS, WARN, FAIL) and actionable remedies.
- **`internal/tui/styles` Package**: Centralized styling using `charmbracelet/lipgloss`.
  - Color palette (Rose Pine inspired).
  - Reusable styles for frames, titles, lists, and status indicators.
- **`xoje doctor` CLI Command**: 
  - Validates Go version (>= 1.22) and Node.
  - Probes Go proxy and GitHub connectivity.
  - Validates `config.json` schema and integrity.
- **Visual Refresh**:
  - CLI output wrapped in frames and status icons.
  - TUI menu with better layout and highlight styles.
- **Main Wiring**: Replace `diagnose` with `doctor` in CLI dispatch and TUI menu.

### Out of Scope
- Skill registry management (browser/preview) — deferred to Phase 3.
- Multi-agent configuration adapters — deferred.
- Advanced update logic (version resolution).

## Capabilities

### New Capabilities
- `health-engine`: Extensible diagnostic framework with remediation tips.
- `branded-ui`: Standardized visual language for the whole tool suite.

### Modified Capabilities
- `environment-diagnosis`: Upgraded to deep health checks (`doctor`).
- `interactive-tui`: Visually refreshed and updated to trigger `doctor`.
- `cli-commands`: `diagnose` removed/aliased to `doctor`.

## Approach

1. **Styling**: Define `internal/tui/styles` first to use in the new output.
2. **Doctor Core**: Implement the `doctor` package with a `Check` interface and a `Runner`.
3. **Check Implementation**:
   - `VersionCheck`: parse `go version` and compare.
   - `NetworkCheck`: simple HTTP GET with timeout.
   - `ConfigCheck`: validate fields in the loaded `Config`.
4. **Refactor**: Replace `verify.DiagnoseAll` with `doctor.Run`.
5. **UI Polish**: Apply Lipgloss styles to `main` dispatch and TUI `View`.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `cmd/xoje/main.go` | High | Dispatch logic and output formatting refactor. |
| `internal/doctor/` | New | All health check logic. |
| `internal/tui/styles/` | New | Visual constants and style functions. |
| `internal/tui/tui.go` | Med | Visual refresh and "diagnose" -> "doctor" rename. |
| `internal/verify/` | Delete | Logic migrated/superseded by `doctor`. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| TUI layout break | Med | Use golden files or manual E2E check on different terminal sizes. |
| Network probe slow | High | Use short timeouts (3s) and concurrent check execution in the Runner. |
| Go version parsing | Low | Use `runtime.Version()` or regex on `go version` output. |

## Rollback Plan

Delete the new `doctor` and `styles` packages. Revert `main.go` to the previous commit (verified v1.0.0).

## Dependencies

- `github.com/charmbracelet/lipgloss` (New)
- `github.com/charmbracelet/bubbletea` (Existing)

## Success Criteria

- [ ] `xoje doctor` reports PASS for Go 1.26.5 and reachable network.
- [ ] Output includes "Remedy" instructions on failures.
- [ ] TUI menu is colored and framed with Lipgloss.
- [ ] All 1.0.0 tests pass or are updated to match 2.0 signatures.
- [ ] `crap4go` score for new `doctor` package stays below 30.
