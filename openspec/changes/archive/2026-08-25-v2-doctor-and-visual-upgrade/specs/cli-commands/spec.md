# Delta for CLI Commands

## RENAMED Requirements

### Requirement: diagnose subcommand → doctor subcommand

(Reason: The comprehensive health command supersedes prerequisite-only diagnosis.)
(Migration: Existing `diagnose` and `bootstrap` invocations remain aliases for `doctor`.)

## MODIFIED Requirements

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