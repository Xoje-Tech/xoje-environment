# health-engine Specification

## Purpose

Provides a modular, extensible framework for executing system health checks and reporting results with actionable remediation steps. Inspired by the `gentle-ai` doctor engine.

## Requirements

### Requirement: Check Interface

The system MUST define a standard `Check` interface that every individual health check implements.

#### Scenario: Interface implementation
- GIVEN a new check (e.g., `DiskSpaceCheck`)
- WHEN the check is registered in the engine
- THEN it provides a unique ID, a human-readable name, and a `Run` method.

### Requirement: Structured Results

Every check MUST return a result containing a status (PASS, WARNING, FAIL), a detailed message, and an optional remedy.

#### Scenario: Check result schema
- GIVEN a check execution
- WHEN it completes
- THEN it returns a `Result` struct with `Status`, `Detail`, and `Remedy` fields.

### Requirement: Concurrent Runner

The engine MUST support running multiple checks concurrently with a timeout per check.

#### Scenario: Parallel execution
- GIVEN a set of 5 health checks
- WHEN the runner starts
- THEN all checks execute in parallel, and results are aggregated once all finish or timeout.

### Requirement: Remediation Action

If a check returns a status other than PASS, it SHOULD include a clear instruction (Remedy) for the user to resolve the issue.

#### Scenario: Remedy reporting
- GIVEN a failing Go version check
- WHEN the report is rendered
- THEN it includes the text: "Install Go 1.22+ via mise: mise use -g go@latest".

## Traceability

- SC-1: `xoje doctor` modular health engine.
- SC-2: Actionable remedies in output.
