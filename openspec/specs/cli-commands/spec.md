# cli-commands Specification

## Purpose

Routes the `xoje` binary's arguments to the correct capability: `diagnose`, `install <tool>`, or `tui` (default when no subcommand is given).

## Requirements

### Requirement: Default subcommand

With no arguments, the CLI MUST route to `tui`.

#### Scenario: No arguments

- GIVEN invocation with no arguments
- WHEN the CLI parses
- THEN the routed command is `tui`

### Requirement: doctor subcommand

The `doctor` subcommand MUST run the complete default diagnostic suite and report status, detail, available remedies, and readiness. The CLI MUST retain `diagnose` and `bootstrap` as aliases.

(Previously: `diagnose` reported only go, node, and git availability, with `bootstrap` as an optional alias.)

#### Scenario: Doctor invocation

- GIVEN `xoje doctor`
- WHEN the CLI routes the command
- THEN the complete default diagnostic suite runs and a doctor report is produced

#### Scenario: Legacy alias

- GIVEN `xoje diagnose` or `xoje bootstrap`
- WHEN the CLI routes the command
- THEN it behaves as `xoje doctor`

### Requirement: install subcommand

The `install` subcommand MUST require a tool name argument and MUST route to the install capability with that tool; a missing tool name MUST produce a usage error.

#### Scenario: install with tool

- GIVEN `xoje install gentle-ai`
- WHEN the CLI routes
- THEN the install capability runs for `gentle-ai`

#### Scenario: install without tool

- GIVEN `xoje install` with no tool argument
- WHEN the CLI parses
- THEN a usage error is returned

### Requirement: Unknown subcommand rejection

An unrecognized subcommand MUST be rejected with an error.

#### Scenario: Unknown subcommand

- GIVEN `xoje frobnicate`
- WHEN the CLI parses
- THEN an error is returned identifying the unknown subcommand

### Requirement: Exit codes

On a routing or execution error the CLI MUST exit non-zero; on success it MUST exit zero.

#### Scenario: Success exit

- GIVEN a successful diagnose or install run
- WHEN the process exits
- THEN the exit code is 0

#### Scenario: Error exit

- GIVEN a parse or execution failure
- WHEN the process exits
- THEN the exit code is non-zero

### Requirement: update subcommand

The `update` subcommand MUST accept an optional tool name and route to the update capability.

#### Scenario: update invocation
- GIVEN `xoje update`
- WHEN the CLI routes
- THEN tool update runs for all registered tools

## Traceability

- SC-2/SC-3/SC-5: routing of diagnose/install/tui
- SC-6: E2E — `go build -o bin/xoje` + subcommands verified (parse + dispatch + exit codes)
