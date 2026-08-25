# Delta for Health Engine

## ADDED Requirements

### Requirement: Bounded concurrent checks

The engine MUST run checks concurrently, preserve one result per check, and fail checks exceeding three seconds with a remedy.

#### Scenario: Concurrent completion

- GIVEN multiple checks complete within three seconds
- WHEN the engine runs them
- THEN every result returns without serial execution

#### Scenario: Check timeout

- GIVEN a check exceeds three seconds
- WHEN the engine runs it
- THEN its result fails and identifies the timeout
- AND includes a retry remedy

### Requirement: Aggregate readiness

The health engine MUST reduce results to `READY`, `DEGRADED`, or `NOT READY`: any failure MUST yield `NOT READY`, otherwise any warning MUST yield `DEGRADED`, and only passing results MUST yield `READY`.

#### Scenario: Warning-only degradation

- GIVEN passing results and at least one warning but no failure
- WHEN readiness is aggregated
- THEN readiness is `DEGRADED`