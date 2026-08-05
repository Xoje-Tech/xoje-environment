# config-state Specification

## Purpose

Manages the persistent xoje configuration: a JSON file holding the active persona and the registry of installed tools. Provides load-or-create semantics so the CLI works on first run.

## Requirements

### Requirement: Load-or-create

The config loader MUST load an existing config file; when the file does not exist, it MUST create the parent directory, write a new config with default values, and return it.

#### Scenario: First run creates defaults

- GIVEN no config file at the path
- WHEN the config is loaded
- THEN a new config is created with the default persona and an empty installed-tools list
- AND the file is persisted to disk

#### Scenario: Existing config loads

- GIVEN a valid config file at the path
- WHEN the config is loaded
- THEN the stored persona and installed-tools values are returned unchanged

### Requirement: JSON format

The config MUST be stored as JSON with snake_case keys `active_persona` and `installed_tools`.

#### Scenario: Persisted shape

- GIVEN a saved config
- WHEN the file is read
- THEN it parses as JSON containing `active_persona` and `installed_tools`

### Requirement: Save round-trip

Changes made to a loaded config (e.g. appending an installed tool) MUST be persisted by Save and MUST be readable on a subsequent load.

#### Scenario: Save then reload

- GIVEN a loaded config with a tool appended
- WHEN Save is called and the config is loaded again
- THEN the reloaded config contains the appended tool

### Requirement: Invalid config handling

Loading a file that exists but is not valid JSON MUST return an error instead of silently overwriting it.

#### Scenario: Corrupt file

- GIVEN a config path whose file contains invalid JSON
- WHEN the config is loaded
- THEN an error is returned and the file is left untouched

### Requirement: Default location

The CLI MUST use `~/.config/xoje/config.json` as the config path when the user does not specify one.

#### Scenario: Default path used

- GIVEN a user with home directory H
- WHEN the CLI starts without an explicit config path
- THEN the config is loaded from `H/.config/xoje/config.json`

## Traceability

- SC-3: install records the tool in config (load-or-create, save round-trip)
