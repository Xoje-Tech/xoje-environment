# Proposal: First-Party Repository-Owned Skills Module

## Intent
Give `xoje-environment` a reproducible, auditable home for first-party skills without conflating repository content with Hermes global skills or `gentle-ai` internals. The outcome is explicit discovery, validation, and controlled projection to verified agent targets.

## Scope

### In Scope
- Canonical `skills/<skill-id>/SKILL.md` content with validated metadata and ownership.
- Manifest/catalog validation, deterministic registry generation, and drift detection.
- Go services with filesystem/process ports for list, verify, target discovery, dry-run, and idempotent sync.
- Preview-based onboarding/migration preserving unmanaged files, detecting collisions, and applying explicit overwrite/backup policy.
- Read-only doctor checks and narrow CLI/TUI wiring after design.

### Out of Scope
- Production implementation in this proposal phase.
- Implicit mutation, vendoring, or overwriting of `~/.hermes/skills/`.
- Symlinks, shell hooks, plugin execution, network installation, or embedded RPI.
- Adapters for agents lacking a verified local contract, target path, and capabilities.
- Replacing `.atl/skill-registry.md`, which remains a generated discovery index.

## Capabilities

### New Capabilities
- `first-party-skills`: canonical content, validation, registry, projection, onboarding, migration, and rollback-safe sync.

### Modified Capabilities
- `cli-commands`: specified skill lifecycle commands only.
- `environment-diagnosis`: read-only skill-integrity checks.
- `config-state`: versioned skill state only if design proves it necessary.
- `interactive-tui`: skill status/onboarding through application ports.

## Approach
Use a Go Onion/Hexagonal module: domain rules for IDs, metadata, ownership, capabilities, and confinement; application services for verify/list/sync; infrastructure adapters for filesystem and allowlisted targets. Sync MUST be confined, atomic, idempotent, non-destructive by default, and dry-run capable. Hermes owns `~/.hermes/skills/`; Xoje may document or delegate native `gentle-ai` operations, never reimplement or silently invoke them as project sync.

Initial agent support is undecided. Qualification requires verified contract, target discovery, write semantics, capability set, and fixture-based adapter tests. Unsupported targets return an explicit safe result.

## Affected Areas
`skills/` and registry (new); `internal/` services/adapters/checks (new or modified); `cmd/xoje/main.go`, `internal/cli`, `internal/tui` (future wiring); `internal/config` (optional state); `.atl/skill-registry.md` (generated only).

## Risks and Rollback
Prevent escape/collisions through canonicalization, confinement, atomic staging, and overwrite policy. Prevent drift with deterministic generation and an allowlist. Rollback disables commands, removes only recorded Xoje-managed targets, restores backups, and never deletes unmanaged or Hermes-global content.

## Dependencies
Existing Go CLI/TUI, config, doctor, and installer boundaries; OpenSpec decisions. No agent capability is assumed.

## Acceptance-Level Outcomes
- [ ] Specs define manifest/version, registry, ownership, collision, backup, rollback, and unsupported-target semantics.
- [ ] Design proves ports/adapters preserve architecture and ownership boundaries.
- [ ] Scenarios cover verify, dry-run, idempotent sync, migration safety, and rollback.
- [ ] Initial agent matrix contains only verified capabilities.

## Unresolved Decisions
Manifest/version and registry format; first agent versus verify/list/sync-only; ownership marker and backup retention; skill state in config versus derived state.
