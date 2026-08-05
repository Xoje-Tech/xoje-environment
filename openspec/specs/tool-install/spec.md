# tool-install Specification

## Purpose

Installs supported developer tools on the host via `go install` and records successful installs in the config registry. Supports dry-run and GOBIN override. Part of the xoje-environment PoC.

## Requirements

### Requirement: Supported tool whitelist

The installer MUST accept only tools registered in its whitelist; any other tool name MUST be rejected with an error before any command execution.

#### Scenario: Known tool accepted

- GIVEN a tool named `gentle-ai`
- WHEN the installer is invoked with that name
- THEN the installer proceeds without a whitelist error

#### Scenario: Unknown tool rejected

- GIVEN a tool name not in the whitelist (e.g. `vim`)
- WHEN the installer is invoked with that name
- THEN an error is returned and no `go install` command is executed

### Requirement: Dry-run mode

In dry-run mode the installer MUST NOT execute `go install`, MUST NOT modify the config registry, and MUST return success for a whitelisted tool.

#### Scenario: Dry-run of a known tool

- GIVEN dry-run enabled and tool `gentle-ai`
- WHEN the installer runs
- THEN it returns success
- AND no install command is executed and no registry entry is written

### Requirement: GOBIN override

When a custom bin directory is supplied, the installer MUST run `go install` with GOBIN set to the absolute path of that directory; when none is supplied, GOBIN MUST be left unset (Go default).

#### Scenario: Custom bin directory

- GIVEN a custom bin directory `tmp/bin`
- WHEN the installer executes `go install`
- THEN the subprocess environment sets GOBIN to the absolute path of `tmp/bin`

#### Scenario: No override

- GIVEN no custom bin directory
- WHEN the installer executes `go install`
- THEN GOBIN is not overridden in the subprocess environment

### Requirement: Install command and failure reporting

For `gentle-ai`, the installer MUST execute `go install github.com/gentleman-programming/gentle-ai/v2/cmd/gentle-ai@latest` and, on failure, MUST return an error that includes the subprocess output.

#### Scenario: Successful install

- GIVEN a host with the Go toolchain and network access
- WHEN the installer runs `go install` for `gentle-ai`
- THEN it returns success

#### Scenario: Failed install

- GIVEN a `go install` that exits non-zero
- WHEN the installer runs
- THEN an error is returned that includes the captured output

### Requirement: Config registry update

After a successful non-dry-run install, the system MUST append the installed tool to the config's installed-tools registry and MUST persist the config.

#### Scenario: Registry records the tool

- GIVEN a successful install of `gentle-ai`
- WHEN the install flow completes
- THEN `gentle-ai` is present in the persisted installed-tools list

## Traceability

- SC-3: `xoje install gentle-ai` installs and records the tool in config (whitelist, GOBIN, registry update)
- SC-4: `gentle-ai version` succeeds post-install (install command correctness)
