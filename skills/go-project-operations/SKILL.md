<!-- xoje-owned:v1 id=go-project-operations -->
---
schema: xoje.skill/v1
id: go-project-operations
category: operations
version: 1.0.0
metadata:
  name: go-project-operations
  summary: Operate and verify the Xoje Go project with reproducible standard-library-first workflows.
capabilities: [documentation]
targets: []
---

# Go Project Operations

Use this project-owned skill when changing `xoje-environment`.

## Workflow

1. Inspect the current branch and working tree before editing.
2. Keep command parsing in `internal/cli` and keep `cmd/xoje` as a thin dispatcher.
3. Prefer the Go standard library and wrap errors with operation context.
4. Write focused table-driven tests before production changes.
5. Run `gofmt`, focused tests, `go test ./...`, and `go vet ./...` before delivery.

Do not modify global Hermes state or user configuration as part of project-local work.
