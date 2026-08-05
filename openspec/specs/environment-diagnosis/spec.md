# environment-diagnosis Specification

## Purpose

Checks whether the host satisfies the prerequisites for tool installation by resolving `go`, `node`, and `git` on the system PATH, and reports aggregate readiness.

## Requirements

### Requirement: Prerequisite set

Diagnosis MUST check exactly the prerequisites `go`, `node`, and `git`.

#### Scenario: Full prerequisite check

- GIVEN a host
- WHEN diagnosis runs
- THEN availability is reported for `go`, `node`, and `git`

### Requirement: LookPath-based availability

For each prerequisite, availability MUST be determined with `exec.LookPath`; the result MUST include the tool name, a boolean availability flag, and the resolved path (empty when unavailable).

#### Scenario: Tool present

- GIVEN `go` resolves on PATH
- WHEN diagnosis checks `go`
- THEN the result reports available with the resolved path

#### Scenario: Tool absent

- GIVEN a name that does not resolve on PATH
- WHEN diagnosis checks it
- THEN the result reports unavailable with an empty path

### Requirement: Aggregate readiness

Diagnosis MUST aggregate the per-prerequisite results and MUST indicate whether all prerequisites are available, identifying any missing ones.

#### Scenario: All prerequisites present

- GIVEN go, node, and git all resolve
- WHEN diagnosis completes
- THEN readiness is reported as fully satisfied

#### Scenario: Prerequisite missing

- GIVEN at least one prerequisite absent
- WHEN diagnosis completes
- THEN readiness is reported as not fully satisfied
- AND the missing tools are identified

## Traceability

- SC-2: `xoje diagnose` reports go/node/git availability
