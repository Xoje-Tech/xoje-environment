# environment-diagnosis Specification (v2.0)

## Purpose

Upgrades the basic prerequisite checks into deep health validations, covering version requirements, network reachability, and configuration integrity.

## Requirements

### Requirement: Prerequisite set

Diagnosis MUST check Go, Node, Git, canonical xoje `PATH`, HTTP access to the Go proxy and GitHub, configuration validity, registered tools, Engram health, and Gentle-AI ecosystem health.

(Previously: Diagnosis checked only whether `go`, `node`, and `git` resolved on `PATH`.)

#### Scenario: Full prerequisite check

- GIVEN a configured host
- WHEN diagnosis runs
- THEN it reports results for runtimes, canonical `PATH`, both network targets, configuration, registered tools, Engram health, and Gentle-AI ecosystem health

### Requirement: Go Version Validation

The system MUST check the Go toolchain version and ensure it is >= 1.22.

#### Scenario: Go version check
- GIVEN Go version 1.26.5 installed
- WHEN `doctor` runs the Go check
- THEN it reports PASS with the detected version.

#### Scenario: Old Go version
- GIVEN Go version 1.21 installed
- WHEN `doctor` runs the Go check
- THEN it reports FAIL and provides the remedy "Upgrade to Go 1.22+".

### Requirement: Connectivity Probes

The system MUST verify that the Go proxy and GitHub are reachable.

#### Scenario: Network check
- GIVEN a host with internet access
- WHEN `doctor` probes `https://proxy.golang.org`
- THEN it reports PASS if the status is 200 OK within 3 seconds.

### Requirement: Config Health Check

The system MUST validate the local `config.json` file for readability and presence of required fields (e.g., `active_persona`).

#### Scenario: Valid config
- GIVEN a valid JSON config at `~/.config/xoje/config.json`
- WHEN `doctor` runs
- THEN it reports PASS for the configuration check.

### Requirement: Aggregate readiness

Diagnosis MUST report `READY` when every check passes, `DEGRADED` when at least one check warns and none fail, and `NOT READY` when any check fails.

(Previously: Diagnosis reported only whether all prerequisites were available and identified missing tools.)

#### Scenario: All checks pass

- GIVEN every diagnostic check passes
- WHEN diagnosis completes
- THEN readiness is `READY`

#### Scenario: Warning without failure

- GIVEN at least one warning and no failed checks
- WHEN diagnosis completes
- THEN readiness is `DEGRADED`

#### Scenario: Failed check

- GIVEN at least one failed check
- WHEN diagnosis completes
- THEN readiness is `NOT READY`

## Traceability

- SC-1: `xoje doctor` reports PASS for Go 1.26.5 and network.
- SC-2: Actionable remedies for version/network issues.
