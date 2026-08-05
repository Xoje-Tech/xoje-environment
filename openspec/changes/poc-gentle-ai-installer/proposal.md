# Proposal: xoje-environment PoC — gentle-ai installer

## Intent

Greenfield Go CLI+TUI to manage/install/sync dev environment across agent hosts. PoC: auto-install `gentle-ai` via `go install`, diagnose host prerequisites, drive from a TUI. Replaces manual, host-specific setup with repeatable tooling.

## Scope

### In Scope
- Go module `github.com/Xoje-Tech/xoje-environment` (Go 1.22+), `.gitignore`
- `internal/config`: JSON load-or-create config (persona, installed tools) — TDD
- `internal/verify`: LookPath checks (go/node/git) — TDD
- `internal/install`: `go install` engine (dry-run, GOBIN override); records tool — TDD
- `internal/cli`: parse `diagnose`/`install`/`tui` (default tui) — TDD
- `internal/tui`: Bubbletea shell (Up/Down/Enter, q/esc/ctrl+c) — TDD
- `cmd/xoje/main.go` wiring + E2E (build, diagnose, install, `gentle-ai version`)

### Out of Scope
- `agents/` adaptor layer; `env/` beyond LookPath checks
- Config sync/distribution across hosts
- `gentle-ai doctor` + RDD/MCP engine sync (deferred; E2E uses `gentle-ai version`)
- Tools other than gentle-ai (engine rejects unknown tools)

## Capabilities

### New Capabilities
- `tool-install`: install via `go install` (dry-run, GOBIN override); updates config registry
- `environment-diagnosis`: LookPath prerequisite checks (go, node, git)
- `config-state`: load-or-create JSON config (persona, installed tools)
- `cli-commands`: subcommand routing — diagnose/install/tui (default)
- `interactive-tui`: Bubbletea shell with keyboard navigation

### Modified Capabilities
None — greenfield; no `openspec/specs/`.

## Approach

Domain-driven packages per plan: `internal/{config,verify,install,cli,tui}` + `cmd/xoje`. Strict TDD (`go test ./...`), stdlib JSON, Bubbletea v1.0.0. `main`: config load → `cli.Parse` → dispatch.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `go.mod`, `go.sum` | New | module + bubbletea |
| `cmd/xoje/main.go` | New | entrypoint wiring |
| `internal/{config,verify,install,cli,tui}/` | New | 5 packages + tests |
| `.gitignore` | New | bin/, dist/, xoje, .env |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Go toolchain missing on host | High | Install Go 1.22+ before apply/verify |
| Network needed for go get/install | Med | Dry-run default in tests |
| Config format drift (YAML vs JSON) | Low | Keep JSON; temp-path tests |

## Rollback Plan

Greenfield — no prod impact. Remove `openspec/changes/poc-gentle-ai-installer/` and generated Go files; `git revert` if repo initialized.

## Dependencies

- Go 1.22+ toolchain (**not installed on host**)
- `github.com/charmbracelet/bubbletea@v1.0.0`
- `github.com/gentleman-programming/gentle-ai/v2` (install target)
- Network access to `proxy.golang.org`

## Success Criteria

- [ ] `go test ./...` passes (5 packages, TDD)
- [ ] `xoje diagnose` reports go/node/git availability
- [ ] `xoje install gentle-ai` installs + records tool in config
- [ ] `gentle-ai version` succeeds post-install
- [ ] `xoje` launches TUI; q/ctrl+c exits
- [ ] E2E: `go build -o bin/xoje` + subcommands verified
