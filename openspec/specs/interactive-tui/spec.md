# interactive-tui Specification

## Purpose

An interactive Bubbletea shell that lets the user choose an action (diagnose, install gentle-ai, exit) with keyboard navigation, and exits cleanly via q/esc/ctrl+c.

## Requirements

### Requirement: TUI launch

The TUI MUST start as a Bubbletea program when the `tui` subcommand (or no subcommand) is selected.

#### Scenario: Launch

- GIVEN `xoje` invoked with no arguments
- WHEN the CLI dispatches
- THEN a Bubbletea TUI session starts

### Requirement: Menu rendering

The TUI view MUST render a header identifying xoje-environment and the selectable choices, with a cursor marker on the currently selected choice.

#### Scenario: Initial view

- GIVEN a freshly created TUI model
- WHEN the view is rendered
- THEN it contains the `xoje-environment` header and the choices with the cursor on the first choice

### Requirement: Cursor navigation

Up/Down arrow keys (and k/j) MUST move the cursor between choices and MUST clamp at the first and last choices.

#### Scenario: Move down and up

- GIVEN the cursor on the first choice
- WHEN Down is pressed twice and Up once
- THEN the cursor ends on the second choice

#### Scenario: Clamp at boundaries

- GIVEN the cursor on the first choice
- WHEN Up is pressed
- THEN the cursor stays on the first choice

### Requirement: Quit keys

q, esc, and ctrl+c MUST quit the TUI program.

#### Scenario: Quit via keys

- GIVEN a running TUI
- WHEN q, esc, or ctrl+c is pressed
- THEN the program exits

### Requirement: Selection behavior

Pressing Enter (or space) on the exit choice MUST quit the program.

#### Scenario: Enter on exit choice

- GIVEN the cursor on the exit choice
- WHEN Enter is pressed
- THEN the program exits

## Traceability

- SC-5: `xoje` launches TUI; q/ctrl+c exits
