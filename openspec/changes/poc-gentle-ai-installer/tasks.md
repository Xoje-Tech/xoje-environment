# Tasks: xoje-environment PoC — gentle-ai installer

## Review Workload Forecast
Decision needed before apply: No (Resolved: Split into 2 batches, size exception handled)
Chained PRs recommended: No
Chain strategy: single local chain (PoC)
400-line budget risk: High (Resolved: Total ~900 lines including remediations)

## Phase 1: Foundation

- [x] 1.1 Create `go.mod` (module `github.com/Xoje-Tech/xoje-environment`, go 1.22), `.gitignore` (`bin/`, `dist/`, `xoje`, `.env`), hello-world `cmd/xoje/main.go`; `git init`; verify `go build ./...`; commit `chore: initial repository structure and hello world`.

## Phase 2: Core packages (TDD)

- [x] 2.1 RED `internal/config/config_test.go`: load-or-create defaults+persist, existing load unchanged, save round-trip, corrupt JSON→error+untouched, `DefaultConfigPath`=`H/.config/xoje/config.json` (t.Setenv HOME); `go test ./internal/config/... -v` FAIL. GREEN `config.go`: `Config{ActivePersona,InstalledTools,path}` snake_case tags, `LoadOrCreate`, `DefaultConfigPath`, `Save` (MarshalIndent); PASS → commit `feat(config): implement JSON-backed config loader and saver`.
- [x] 2.2 RED `internal/verify/verify_test.go`: present (stub bin on PATH)/absent, aggregate all-present + missing-identified; `go test ./internal/verify/... -v` FAIL. GREEN `verify.go`: `DiagnosticResult{Name,Path,Available}`, `VerifyExecutable` (LookPath), `DiagnoseAll` (go/node/git + readiness); PASS → commit `feat(verify): implement LookPath-based diagnostic check`.
- [x] 2.3 RED `internal/install/install_test.go` (stub `go` script on PATH via t.Setenv, clear GOBIN): whitelist accept, unknown `vim`→error+no exec, dry-run→nil+no exec, GOBIN=abs(tmp/bin), no override→unset, failure→error contains output; `go test ./internal/install/... -v` FAIL. GREEN `install.go`: whitelist `gentle-ai`→`github.com/gentleman-programming/gentle-ai/v2/cmd/gentle-ai@latest`; `InstallTool(name,customBinDir,dryRun)`: gate→dry-run→GOBIN (filepath.Abs) or unset→CombinedOutput→error wraps output; PASS → commit `feat(install): add engine to execute tool installation via go install`.

## Phase 3: CLI + TUI (TDD)

- [x] 3.1 RED `internal/cli/cli_test.go` table: no args→tui, diagnose, bootstrap→diagnose, install+tool, install w/o tool→usage error, unknown→error; `go test ./internal/cli/... -v` FAIL. GREEN `cli.go`: `Parse(args)(cmd,tool,err) default tui, alias, required tool, unknown rejected; PASS → commit `feat(cli): add lightweight CLI subcommand parser`.
- [x] 3.2 RED `internal/tui/tui_test.go` (after `go get github.com/charmbracelet/bubbletea@v1.0.0`): header+choices+cursor first, down×2/up→second, clamp both ends, q/esc/ctrl+c→Quit, enter/space on exit→Quit; `go test ./internal/tui/... -v` FAIL. GREEN `tui.go`: `Model{cursor,choices}` "diagnose"/"install gentle-ai"/"exit", `NewModel`, `Init`, `Update` (up/k down/j clamp, quit, enter-on-exit), `View` (cursor `>`); PASS → commit `feat(tui): setup Bubbletea state machine and primary rendering engine`.

## Phase 4: Wiring

- [x] 4.1 Rewrite `cmd/xoje/main.go`: `config.LoadOrCreate(config.DefaultConfigPath())` → `cli.Parse(os.Args[1:])` → diagnose (`verify.DiagnoseAll` report), install (`InstallTool(tool,"",false)`; success→append+Save per D5), tui (`tea.NewProgram(tui.NewModel()).Run()`); exit 0/1. `go test ./...` PASS, `go vet ./...`, `gofmt -l .` empty → commit `feat(main): integrate config, cli, verify and tui into main executable`.

## Phase 5: E2E (needs Go 1.22+ toolchain + network)

- [x] 5.1 Provision Go 1.26.5 (absent on host) and network to proxy.golang.org.
- [x] 5.2 `go build -o bin/xoje cmd/xoje/main.go`; `./bin/xoje diagnose` reports go/node/git; `./bin/xoje install gentle-ai`; `gentle-ai version`; `~/.config/xoje/config.json` has gentle-ai; `./bin/xoje` TUI q/ctrl+c exits; `./bin/xoje frobnicate` exits non-zero.
- [x] 5.3 Final gate `go test ./...` + `go build ./...`; capture evidence for verify-report.

## Phase 6: Quality Remediation (Post-Verify)

- [x] 6.1 Refactor `main.go` to use a testable `Run()` function.
- [x] 6.2 Add `cmd/xoje/main_test.go` with coverage for `Run()`.
- [x] 6.3 Verify quality with `crap4go` (main score dropped from 306 to 12).

## Phase 7: Update Command

- [ ] 7.1 RED `internal/cli/cli_test.go`: add `update` cases (no arg, single tool). GREEN `cli.go`: implement `update` routing.
- [ ] 7.2 Implement `update` logic in `main.go`: if tool given, run `install`; if no tool, loop over `cfg.InstalledTools`.
- [ ] 7.3 TDD `main_test.go`: cover `update` flow.
- [ ] 7.4 Add `update` option to TUI choices in `internal/tui/tui.go`.

## Coverage
config-state→2.1; environment-diagnosis→2.2; tool-install→2.3+4.1 (registry); cli-commands→3.1+4.1; interactive-tui→3.2+4.1; threat matrix (subprocess)→2.3; design file table→all; E2E→Phase 5; main-wiring→6.1+6.2.
