# branded-ui Specification

## Purpose

Establishes a professional visual identity for `xoje-environment` using the `charmbracelet/lipgloss` framework. Standards are inspired by the *Rose Pine* color palette used in `gentle-ai`.

## Requirements

### Requirement: Centralized Styling

Visual styles (colors, borders, padding) MUST be centralized in a dedicated package to ensure consistency across CLI and TUI.

#### Scenario: Style reuse
- GIVEN the `internal/tui/styles` package
- WHEN a component needs a title style
- THEN it imports and applies `styles.TitleStyle`.

### Requirement: Color Palette (Rose Pine)

The system MUST use a standardized set of hex colors for different UI states.

#### Scenario: Color definitions
- GIVEN the style definitions
- WHEN hex values are assigned
- THEN they match Rose Pine (e.g., Lavender for titles, Green for success, Red for errors).

### Requirement: Framed Output

Subcommand results (like `doctor` or `install`) MUST be wrapped in visual frames using Lipgloss border styles.

#### Scenario: Result framing
- GIVEN a doctor report
- WHEN rendered to the terminal
- THEN the entire report block is enclosed in a double-line border with a specific border color.

### Requirement: Semantic Status Indicators

Status reports MUST use both color and icons for accessibility and clarity.

#### Scenario: Status icons
- GIVEN a health check result
- WHEN displayed
- THEN PASS shows a green "✅", WARNING shows a yellow "⚠️", and FAIL shows a red "❌".

### Requirement: Doctor report presentation

Doctor reports MUST visually distinguish passing, warning, and failed checks, MUST show details and available remedies, and MUST end with aggregate readiness.

#### Scenario: Mixed doctor results

- GIVEN doctor results containing pass, warning, and failure states with remedies
- WHEN the report is rendered
- THEN each state has a distinct semantic indicator
- AND details, remedies, and `NOT READY` aggregate readiness are visible

## Traceability

- SC-3: TUI menu and CLI output framed with Lipgloss.
- SC-2: Rich output formatting.
