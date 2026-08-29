<!-- xoje-owned:v1 id=openspec-change-conventions -->
---
schema: xoje.skill/v1
id: openspec-change-conventions
category: specification
version: 1.0.0
metadata:
  name: openspec-change-conventions
  summary: Keep OpenSpec changes scoped, testable, and separate from production implementation.
capabilities: [documentation]
targets: []
---

# OpenSpec Change Conventions

Use this project-owned skill when preparing a change proposal or implementation boundary.

## Boundaries

- Record intent and requirements in OpenSpec artifacts, but do not treat them as executable code.
- Keep production behavior in the Go packages under `internal/` and `cmd/`.
- State explicit non-goals, especially network access, global Hermes mutation, and implicit synchronization.
- Turn each accepted behavior into a focused test before implementation.
- Verify the repository after implementation without silently repairing derived artifacts.
