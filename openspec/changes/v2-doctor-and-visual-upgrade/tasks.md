# Tasks: v2.0 — Doctor & Visual Upgrade

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 350-500 |
| 400-line budget risk | Medium |
| Chained PRs recommended | Yes |
| Suggested split | PR 1 (Styles & Engine) → PR 2 (Checks & Wiring) |
| Delivery strategy | ask-on-risk |
| Chain strategy | stacked-to-main |

Decision needed before apply: Yes
Chained PRs recommended: Yes
Chain strategy: stacked-to-main
400-line budget risk: Medium

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|------|------|-----------|----------------------|-----------------|-------------------|
| 1 | Foundation: Styles & Doctor Engine | PR 1 | `go test ./internal/doctor/...` | N/A (Internal API) | `internal/tui/styles`, `internal/doctor/doctor.go` |
| 2 | Implementation: Checks & UI Wiring | PR 2 | `go run ./cmd/xoje doctor` | Manual TUI verification | `internal/doctor/checks.go`, `internal/tui/tui.go`, `cmd/xoje/main.go` |

## Phase 1: Foundation (Styles & Engine)

- [x] 1.1 Create `internal/tui/styles/styles.go` with Rose Pine constants and Lipgloss styles.
- [x] 1.2 RED: Write unit test for `doctor.Runner` with mock checks in `internal/doctor/doctor_test.go`.
- [x] 1.3 Create `internal/doctor/doctor.go` with `Check` interface, `Result` struct, and concurrent `Runner`.
- [x] 1.4 GREEN: Implement `Runner.Run()` logic and satisfy tests.

## Phase 2: Core Implementation (Checks)

- [x] 2.1 RED: Write unit tests for concrete checks in `internal/doctor/checks_test.go`.
- [x] 2.2 Create `internal/doctor/checks.go` implementing `GoVersionCheck`, `NetworkCheck`, and `ConfigCheck`.
- [x] 2.3 GREEN: Implement check logic ensuring correct status/remedy mapping.

## Phase 3: Integration & UI Wiring

- [x] 3.1 Modify `internal/tui/tui.go` to import `styles` and apply themes to components.
- [x] 3.2 Update TUI message handling to use `doctor.Run()` for the diagnosis action.
- [x] 3.3 Modify `cmd/xoje/main.go` to refactor `dispatch` and apply styles to CLI output.

## Phase 4: Testing & Verification

- [x] 4.1 Run full test suite: `go test ./...`.
- [x] 4.2 Verify TUI appearance and concurrent doctor execution manually.
- [x] 4.3 Verify `xoje doctor` CLI command output and framing.

## Phase 5: Cleanup

- [x] 5.1 Remove temporary debug prints (none found).
- [x] 5.2 Archive the SDD change.
