# Tasks: First-Party Skills Module

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | 1,400–2,000 |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Suggested split | Contracts → planning → transactions → integration |
| Delivery strategy | ask-on-risk |
| Chain strategy | feature-branch-chain |

Decision needed before apply: Resolved
Chained PRs recommended: Yes
Chain strategy: feature-branch-chain
400-line budget risk: High

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|---|---|---|---|---|---|
| 1 | Validation and registries | PR 1 | `go test ./internal/skills/... -run 'Manifest|Registry'` | `go run ./cmd/xoje skill verify` | domain/source/registry and `skills/**` |
| 2 | Planning and confinement | PR 2 | `go test ./internal/skills/... -run 'Discover|Plan|Confinement'` | unsupported-target sync dry-run | target/planner adapters |
| 3 | Transactions and recovery | PR 3 | `go test ./internal/skills/... -run 'Transaction|Recover|Rollback'` | temp-root fault injection | transaction adapter/state root |
| 4 | CLI/TUI/doctor integration | PR 4 | `go test ./internal/cli ./internal/tui ./internal/doctor ./cmd/xoje` | no-args TUI; status; doctor | routing/presentation files |

## Phase 1: Contracts and Deterministic Registry

- [ ] 1.1 RED: test `internal/skills/domain/manifest_test.go` for malformed markers, unsafe/duplicate IDs, unknown schema data, and shell-looking Markdown remaining inert data.
- [ ] 1.2 GREEN/REFACTOR: implement typed manifests, ownership, validation, and ports in `internal/skills/domain/**`; keep bodies inert.
- [ ] 1.3 RED: add deterministic/drift cases in `internal/skills/registry/registry_test.go`, including normalized paths, bytewise ordering, fingerprints, LF/final-newline output, and no silent repair.
- [ ] 1.4 GREEN/REFACTOR: implement source loading and registry generation in `internal/skills/{source,registry}/**`; add canonical `skills/**`.

## Phase 2: Safe Discovery and Planning

- [ ] 2.1 RED: test absent/unqualified targets, dry-run/idempotence, collisions, explicit overwrite, path/symlink escapes, and zero writes in `internal/skills/application/planner_test.go`.
- [ ] 2.2 GREEN/REFACTOR: implement list/status/verify/discover/generate and plan-only sync/onboard services in `internal/skills/application/**` with an empty qualified-adapter registry.
- [ ] 2.3 RED→GREEN: implement confined filesystem adapters with `t.TempDir()`; reject ambiguity, network, subprocesses, Hermes roots, and hooks/plugins.

## Phase 3: Transactions, Backup, and Recovery

- [ ] 3.1 RED: test all transitions, lock contention, unfinished recovery, rename failure per index, exact compensation, and `partial_recovery` in `internal/skills/transaction/transaction_test.go`.
- [ ] 3.2 GREEN/REFACTOR: implement typed journals/receipts, lock, staging, ordered renames, reverse compensation, and hash verification in `internal/skills/transaction/**`.
- [ ] 3.3 RED→GREEN: test and implement `168h` retention, finalized GC, retention `0`, retained partial recovery, linked rollback transactions, and compensation-only recover.

## Phase 4: Compatible Integration

- [ ] 4.1 RED→GREEN: extend `internal/cli/cli_test.go` and parser for all specified `skill` commands/flags while preserving legacy routes/exits.
- [ ] 4.2 RED→GREEN: wire services in `cmd/xoje/main.go`; add command characterization tests proving dry-run/read-only behavior and explicit mutation policy.
- [ ] 4.3 RED→GREEN: inject application ports into `internal/tui/tui.go`; test status/onboarding state transitions in `internal/tui/tui_test.go` without filesystem policy in the model.
- [ ] 4.4 RED→GREEN: add read-only skill integrity/recovery checks to `internal/doctor/checks.go` and tests preserving readiness aggregation and no network/write behavior.
- [ ] 4.5 REFACTOR/verify: run `gofmt`, focused tests, `go test ./...`, `go test -cover ./...`, `go vet ./...`, and offline CLI/TUI harnesses; record each work-unit rollback evidence.
