# Technical Design: First-Party Skills Module

## Technical Approach

Add Hexagonal `internal/skills`: CLI/TUI/doctor → application → domain ← source, registry, target, transaction, and backup adapters. `skills/<id>/SKILL.md` is canonical; registries are derived. Hermes global skills, network, subprocesses, symlinks, hooks/plugins, RPI execution, and implicit `gentle-ai` operations remain excluded.

## Decisions

| Decision | Contract and rationale |
|---|---|
| Source path | The approved proposal/spec makes `skills/<id>/SKILL.md` a change-specific project exception to generic `.agents/skills/<category>/<name>/SKILL.md`; `.agents/**` is not generated. Future qualified adapters may render managed copies only under their approved contract. |
| Targets | Initial qualified-adapter matrix is empty. Qualification requires confined roots, marker/layout rules, fixtures, and verified compatibility; therefore no agent support is invented. |
| Transaction | One UUID owns one target, advisory lock, journal, staging, and backups. Sequential atomic renames plus reverse compensation are portable and truthful; multi-file atomicity is not claimed. |

## Contracts

Every canonical file starts exactly with `<!-- xoje-owned:v1 id=<id> -->`, followed by YAML frontmatter:

```yaml
schema: xoje.skill/v1
id: string               # lowercase [a-z0-9-]+, path-safe, globally unique
category: string         # lowercase [a-z0-9-]+
version: string          # SemVer
metadata: {name: string, summary: string}
capabilities: [documentation|workflow|rpi-reference]
targets: [string]        # sorted qualified adapter IDs; initially []
```

Unknown fields, versions, capabilities, targets, malformed markers, or marker/frontmatter ID mismatch fail closed. Managed target files start `<!-- xoje-managed:v1 skill=<id> target=<target> source_sha256=<hex> -->`; mutation also requires the registry’s prior target hash to match.

`.atl/skill-registry.json` is UTF-8/LF/final-newline canonical JSON. Types are explicit:

```text
Registry={schema:"xoje.skill-registry/v1",skills:Skill[]}
Skill={id:string,category:string,version:semver,path:string,
 metadata:{name:string,summary:string},capabilities:Capability[],
 targets:string[],source_sha256:lowercase-hex64}
```

Arrays sort bytewise by ID/value; paths are normalized repository-relative `/` paths. Markdown renders this array only. Generation never repairs drift silently.

`Journal={schema:"xoje.skill-transaction/v1",uuid:UUID,target:string,kind:"apply"|"rollback"|"recover",parent_uuid:UUID|null,state:State,created_at:RFC3339,updated_at:RFC3339,retention_seconds:uint64,plan_sha256:hex64,completed_index:uint32,operations:Operation[]}`; each operation has typed relative `path`, `action:"create"|"replace"|"delete"`, typed pre/post `{exists:bool,sha256:hex64|null,marker:string|null}`, and `backup_path:string|null`. `Receipt={schema:"xoje.skill-receipt/v1",uuid:UUID,parent_uuid:UUID|null,target:string,outcome:State,started_at:RFC3339,finished_at:RFC3339,plan_sha256:hex64,journal_sha256:hex64,pre_fingerprint:hex64,post_fingerprint:hex64}` survives backup deletion.

States are `prepared|applying|rollbacking|rolled_back|partial_recovery|committed|finalized`. Legal transitions: `prepared→applying|rollbacking`; `applying→committed|rollbacking`; `committed→finalized|rollbacking`; `rollbacking→rolled_back|partial_recovery`; a later `recover` reacquires the lock and permits `partial_recovery→rollbacking`. Lock-terminal outcomes are `rolled_back|partial_recovery|finalized`; only the middle one remains recoverable. No other transition is valid.

## CLI and Recovery

No arguments enters TUI; legacy routes remain compatible. `list|status|verify|discover` are read-only; `generate` updates only registries. `sync --target <id> [--dry-run] [--overwrite-unmanaged] [--backup-retention <duration>]`; `onboard --target <id> [--apply]` shares flags. Default retention is `168h`; dry-run creates no lock/artifact. `rollback --transaction <uuid> [--dry-run]` requires current committed hashes and creates a new `kind=rollback` transaction whose `parent_uuid` references the unchanged finalized receipt.

`recover --target <id> [--dry-run]` accepts journals in `prepared|applying|rollbacking|partial_recovery`, only compensates toward recorded pre-state, and never resumes writes. Exact restoration yields `rolled_back`; divergence remains `partial_recovery` with manual paths/hashes.

A non-blocking OS advisory lock on `$XDG_STATE_HOME/xoje/skill-transactions/<target>/lock` is exclusive cross-process; contention exits unchanged with `target_busy`. PID/mtime never proves staleness: kernel lock ownership is authoritative, while a free lock plus unfinished journal triggers recovery. Release occurs only after a terminal state persists.

## Transaction Sequence

```mermaid
sequenceDiagram
 participant C as CLI/TUI
 participant A as Application
 participant J as Journal/Lock
 participant F as Target FS
 C->>A: sync/onboard --apply
 A->>J: acquire lock; persist prepared
 A->>F: preflight, backup, stage
 loop ordered operations
  A->>F: atomic rename
  A->>J: persist applying + index
 end
 alt hashes match
  A->>J: committed then finalized
 else failure
  A->>J: rollbacking
  A->>F: reverse restore and verify
  alt exact
   A->>J: rolled_back
  else divergent
   A->>J: partial_recovery
  end
 end
 A->>J: release only after terminal persist
```

## Backup Lifecycle

Backups are created after locked preflight and before any rename; the transaction exclusively owns its directory. Verified automatic compensation or explicit rollback deletes backup bytes after `rolled_back`. Failed compensation preserves everything in `partial_recovery`. `finalized` backups expire at `finished_at + retention`; the next mutating command performs GC. Retention `0` deletes immediately after `finalized`, never earlier. Receipts are retained independently.

## Files and Tests

`internal/skills/**` adds domain/use cases/adapters/tests; CLI, TUI, and doctor gain compatible wiring; `skills/**`, registries, and transaction state are added. Existing install/config/verify boundaries remain unchanged. Strict TDD covers schemas, deterministic drift, unsupported no-write, confinement, dry-run, consent, lock contention, every rename failure, state transitions, retention/GC, rollback/recover, and read-only doctor behavior.

## Threat Matrix and Rollout

Documentation-like bodies are applicable and remain inert data (RED: shell-looking Markdown never executes). Git selection/commit/push and PR composition are N/A because no VCS or subprocess boundary exists. No schema/data migration is required. Adoption is preview-first onboarding preserving unmanaged files; rollout enables verify/generate first and keeps sync no-write until an adapter qualifies. No blocking open questions remain.
