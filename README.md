# xoje-environment

> Platform-agnostic CLI + TUI tool to manage, install, and sync our development environment across agent hosts.

Built in Go with a focus on portability and zero external dependencies (except the Bubbletea TUI framework). This tool is the **Single Source of Truth (SSOT)** for your development identity and skills.

## 🚀 PoC v1.0.0: gentle-ai Installer

The current version focuses on the automated installation and management of `gentle-ai`, serving as a proof of concept for the environment orchestration engine.

### Prerequisites

- **Go 1.22+** (Recommended: install via [mise](https://mise.jdx.dev/): `mise use -g go@latest`).
- Network access (to download tools via `go install`).

### Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/Xoje-Tech/xoje-environment.git
cd xoje-environment
go build -o bin/xoje cmd/xoje/main.go
```

To install it globally in your `$GOPATH/bin`:

```bash
go install ./cmd/xoje
```

## 🛠 Usage

`xoje` supports both direct CLI commands and an interactive TUI.

### CLI Commands

| Command | Description |
|---------|-------------|
| `xoje diagnose` | Check if prerequisites (`go`, `node`, `git`) are installed on the host. |
| `xoje install gentle-ai` | Install the latest version of `gentle-ai` and record it in config. |
| `xoje tui` | Launch the interactive terminal interface (Default). |

### Interactive TUI

Launch with `xoje` (no arguments).
- **Navigation**: Use Arrow keys or `k`/`j`.
- **Selection**: Press `Enter` to run the selected action.
- **Exit**: Press `q`, `esc`, or `ctrl+c`.

## ⚙️ Configuration

Configuration is stored in a portable JSON format at:
`~/.config/xoje/config.json`

It tracks your `active_persona` and the list of `installed_tools` managed by the environment.

## 🏗 Architecture & Quality

This project follows **Spec-Driven Development (SDD)** and strict **TDD**.

- **Internal Packages**: Decoupled domain packages under `internal/`.
- **Quality Gates**: Audited with Uncle Bob's `crap4go` and `dry4go`.
- **Test Suite**: `go test ./... -v` (includes unit and lifecycle tests).

## 📜 License

Apache-2.0
