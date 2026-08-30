# Delta Specification: First-Party Skills Module

## ADDED Requirements

### Requirement: Canonical source and versions

`skills/<skill-id>/SKILL.md` MUST be canonical. IDs MUST be unique and path-safe. Versioned manifest/frontmatter MUST declare ownership, metadata, capabilities, and targets; unsupported versions MUST fail closed.

#### Scenario: Invalid skill
- GIVEN malformed metadata, duplicate ID, or unsupported version
- WHEN verification runs
- THEN it reports skill/field/remediation and writes nothing

### Requirement: Deterministic registries

The module MUST generate `.atl/skill-registry.json` and `.atl/skill-registry.md` from validated source with version, normalized paths, ordering, and fingerprints. Markdown is an index. Drift MUST fail verify/doctor; commands MUST NOT repair silently.

#### Scenario: Offline generation
- GIVEN unchanged source
- WHEN generation runs twice offline
- THEN outputs are identical and edits are reported as drift

### Requirement: Ownership and capabilities

Xoje MUST distinguish canonical, managed, and unmanaged files and MUST NOT mutate `~/.hermes/skills/`. Only verified allowlisted adapters MAY run; hooks/plugins MUST no-write.

#### Scenario: Unsupported target
- GIVEN an absent, unqualified, or incompatible agent
- WHEN sync is requested
- THEN it reports the reason and mutates nothing

### Requirement: Discovery and sync

Discovery MUST deterministically report present, absent, and unsupported targets. Sync MUST support dry-run, results, offline execution, and idempotence; it MUST NOT invoke network or native commands.

#### Scenario: Dry-run and repeat
- GIVEN a changed skill and supported target
- WHEN dry-run, then unchanged sync, runs
- THEN the first writes nothing and the second is a no-op

### Requirement: Confinement and atomicity

Paths MUST be canonicalized and confined to roots; traversal, absolute escape, symlink escape, and ambiguity MUST fail. Writes MUST stage and atomically replace; failure MUST preserve committed state.

#### Scenario: Unsafe path or failure
- GIVEN an escaping path or failed write
- WHEN verification or sync runs
- THEN it errors, cleans staging, and performs no mutation

### Requirement: Collision, migration, backup, rollback

Sync MUST be non-destructive. Unmarked files MUST be preserved. Overwrite requires explicit policy, backup, and a recoverable record. Onboarding MUST preview/resume. Rollback MUST affect only recorded Xoje paths.

#### Scenario: Collision
- GIVEN an unmarked target
- WHEN onboarding or sync evaluates it
- THEN it skips and requests policy

### Requirement: CLI, TUI, doctor, and state

CLI MUST retain routes while adding list/status, verify, generate, discover, sync, onboarding, and rollback with compatible exits. TUI MUST use application ports. Doctor MUST add read-only checks. State MUST be versioned; derived state is preferred.

#### Scenario: Diagnosis
- GIVEN invalid metadata, drift, or damage
- WHEN `xoje doctor` runs
- THEN it reports remediation without writes or network

### Requirement: RPI/OpenSpec boundary

Skills MAY document RPI/OpenSpec, but Xoje MUST NOT embed an RPI runner, assume proprietary protocols, or claim untested adapters. OpenSpec remains proposal/design/tasks. gentle-ai native sync is separate, never implicit.

#### Scenario: RPI reference
- GIVEN an RPI-described skill
- WHEN offline verification runs
- THEN documentation is accepted without executing a runner

## MODIFIED Requirements

### Requirement: CLI command routing

Existing routing and errors MUST remain compatible as skill routes are added. (Previously: tools and diagnostics only.)

#### Scenario: Legacy command
- GIVEN an existing command or alias
- WHEN parsed
- THEN its route and exit semantics remain unchanged

### Requirement: Environment diagnosis

Diagnosis MUST retain prerequisite and READY/DEGRADED/NOT READY aggregation and add read-only skill checks. (Previously: no skill checks.)

#### Scenario: Healthy module
- GIVEN valid source, matching registries, and intact targets
- WHEN diagnosis completes
- THEN checks pass and aggregation applies

### Requirement: Interactive TUI integration

TUI MAY add skill actions, but MUST keep filesystem policy in application/domain ports. (Previously: doctor and tools only.)

#### Scenario: TUI conflict
- GIVEN a sync conflict
- WHEN displayed
- THEN TUI presents the result without policy

## Residual Design Decisions

Design MUST fix schemas, versions, agent matrix, flags, marker, backup retention, journal, registry formats, and state representation.