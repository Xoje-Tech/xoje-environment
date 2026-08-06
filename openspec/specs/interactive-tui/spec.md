# interactive-tui Specification (v2.0)

## Purpose

Refreshes the TUI with professional styling and updates the "diagnose" action to trigger the new comprehensive "doctor" workflow.

## Requirements

### Requirement: Branded TUI Layout

The TUI MUST apply the `branded-ui` standards, wrapping the entire menu in a Lipgloss frame.

#### Scenario: Framed menu
- GIVEN the TUI is launched
- WHEN rendered
- THEN it displays a double-border frame around the options with a Lavender border.

### Requirement: Enhanced Navigation Feedback

The currently selected item MUST be highlighted using a distinct color and a persistent cursor icon.

#### Scenario: Selection highlight
- GIVEN the menu
- WHEN a choice is focused
- THEN it is displayed in Lavender/Bold text with a prefix icon (e.g., "▸ ").

### Requirement: Command Integration (Doctor)

Selecting the "doctor" (previously "diagnose") option MUST trigger the `health-engine` runner.

#### Scenario: Trigger doctor
- GIVEN the TUI menu
- WHEN "doctor" is selected and Enter is pressed
- THEN the TUI closes and the full framed health report is printed to the terminal.

## Traceability

- SC-3: TUI menu colored and framed with Lipgloss.
- SC-4: 1.0.0 tests updated to match 2.0 (Doctor rename).
