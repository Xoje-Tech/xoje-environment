# Tasks: v2-doctor-and-visual-upgrade

## Implementation
- [x] Create Rose Pine style definitions in `internal/tui/styles`
- [x] Refactor `install.go` to support idempotency and canonical paths
- [x] Implement `doctor` health engine core
- [x] Add PATH validation and binary collision detection in `doctor`
- [x] Integrate `doctor` results into TUI View
- [x] Implement interactive "Remedy" execution from TUI
- [x] Add integration tests for shell PATH edge cases

## Verification
- [x] Unit tests for `internal/install`
- [x] Unit tests for `internal/doctor`
- [x] Unit tests for `internal/tui`
- [x] Manual verification of binary cleanup

## Continuation Evidence (2026-08-23)
- `go test ./cmd/xoje -run '^TestDoctorReportIncludesAggregateReadiness$' -count=1` — exit 0.
- `go test ./internal/tui -run '^TestDoctorRemedySelection$' -count=1` — exit 0.
- `go test ./internal/doctor -run '^TestEnvPathCheckCanonicalBinDir$' -count=1` — exit 0.
- `go test ./internal/cli ./cmd/xoje ./internal/doctor ./internal/install ./internal/tui -count=1` — exit 0; all 5 focused packages passed.
- `go test ./... -count=1` — exit 0; all 7 packages with tests passed, and `internal/tui/styles` reported no test files.
- Safe runtime CLI harness: built `./cmd/xoje` into `/tmp/xoje-v2-step3-runtime.GOSZ6P`, ran `xoje version` with isolated `HOME`/`XDG_CONFIG_HOME`, and received `xoje version 1.1.0` with exit 0. The real config checksum remained `8f8ea880e611dc664849cf772083624ec445bf5d49159dd03d57e83b1b1125b5` before and after.

## Final Local Evidence (2026-08-24)
- Implementation checklist: 11/11 tasks complete.
- `gofmt -l .` — no output.
- `go test ./...` — exit 0.
- `go test -race ./...` — exit 0.
- Total test coverage — 40.2%.
- `go vet ./...` — exit 0.
- `go build ./...` — exit 0.
- `git diff --check` — exit 0.
- Dev tracker: WU1-WU3, Design Phase, and Task Phase are Done; the milestone remains open until delivery.

## Formal SDD Status
- Proposal, design, and tasks are complete, with 11/11 task checkboxes complete.
- Specs, apply-progress, and verify-report are absent.
- Apply, verify, and archive remain formally blocked; the next recommended phase is spec.
- The passing local evidence above does not constitute formal SDD verification or archival.

## Accepted Scope Deviation
- `internal/install/install.go` contains advanced updater version-resolution logic even though `proposal.md` declares that logic out of scope.
- This is an accepted and documented implementation risk, not a requirement satisfied by this change.
