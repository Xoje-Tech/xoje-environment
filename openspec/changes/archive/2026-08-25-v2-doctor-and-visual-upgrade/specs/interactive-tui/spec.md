# Delta for Interactive TUI

## ADDED Requirements

### Requirement: Interactive doctor remedies

The TUI MUST expose a doctor view, navigate its results, and dispatch an available remedy action. Selecting a result without an action MUST NOT exit the view.

#### Scenario: Apply available remedy

- GIVEN the doctor view contains a selected result with a remedy action
- WHEN Enter is pressed
- THEN the TUI exits with that action selected for dispatch

#### Scenario: Result has no action

- GIVEN the doctor view contains a selected result without a remedy action
- WHEN Enter is pressed
- THEN the TUI remains in the doctor view