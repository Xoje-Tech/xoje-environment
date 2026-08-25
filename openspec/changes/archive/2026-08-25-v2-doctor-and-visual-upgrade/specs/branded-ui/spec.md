# Delta for Branded UI

## ADDED Requirements

### Requirement: Doctor report presentation

Doctor reports MUST visually distinguish passing, warning, and failed checks, MUST show details and available remedies, and MUST end with aggregate readiness.

#### Scenario: Mixed doctor results

- GIVEN doctor results containing pass, warning, and failure states with remedies
- WHEN the report is rendered
- THEN each state has a distinct semantic indicator
- AND details, remedies, and `NOT READY` aggregate readiness are visible