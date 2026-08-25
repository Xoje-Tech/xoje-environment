# Delta for Environment Diagnosis

## MODIFIED Requirements

### Requirement: Prerequisite set

Diagnosis MUST check Go, Node, Git, canonical xoje `PATH`, HTTP access to the Go proxy and GitHub, configuration validity, registered tools, Engram health, and Gentle-AI ecosystem health.

(Previously: Diagnosis checked only whether `go`, `node`, and `git` resolved on `PATH`.)

#### Scenario: Full prerequisite check

- GIVEN a configured host
- WHEN diagnosis runs
- THEN it reports results for runtimes, canonical `PATH`, both network targets, configuration, registered tools, Engram health, and Gentle-AI ecosystem health

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