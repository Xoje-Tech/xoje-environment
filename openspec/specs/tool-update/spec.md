# tool-update Specification

## Purpose

Updates installed developer tools to their latest versions by re-running the installation logic. Supports updating a single tool or the entire registered fleet.

## Requirements

### Requirement: Single tool update

The `update` command MUST accept an optional tool name and execute the installation logic for it.

#### Scenario: Update specific tool
- GIVEN tool `gentle-ai` is already in the config registry
- WHEN running `xoje update gentle-ai`
- THEN the installer executes the update for `gentle-ai`

### Requirement: Fleet update (all tools)

If no tool name is provided, the system MUST iterate through all tools in `installed_tools` and update each one.

#### Scenario: Update all tools
- GIVEN `gentle-ai` and another tool are in the config
- WHEN running `xoje update`
- THEN the system executes the update logic for every tool in the registry

### Requirement: Error handling

If a tool is not in the whitelist or not in the registry, it MUST return an error.

## Traceability

- SC-7: `xoje update` refreshes installed tools from the registry.
