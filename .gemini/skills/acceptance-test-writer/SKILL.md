---
name: acceptance-test-writer
description: Generates Acceptance tests using BDD principles. These tests use Gherkin (Given/When/Then) and the `godog` framework to verify that the application satisfies high-level user requirements in a readable format.
---

# Acceptance Test Writer

## Overview
Acceptance tests bridge the gap between user requirements and technical implementation. They use the Gherkin language to describe behaviors in plain English, which are then mapped to Go code using `godog`.

## Workflow
1. **Define Feature**: Create a `.feature` file in `tests/acceptance/features/`.
2. **Write Scenarios**: Use Gherkin syntax:
    - `Given`: Initial state.
    - `When`: User action.
    - `Then`: Expected outcome.
3. **Map Steps**: Implement Go functions that correspond to the Gherkin steps.
4. **Setup Godog Suite**: Initialize the test suite and context in `tests/acceptance/acceptance_test.go`.
5. **Verify**: Run `go test ./tests/acceptance/...`

## Patterns

### Gherkin Feature Template
```gherkin
Feature: [Feature Name]
  Scenario: [Scenario Name]
    Given [state]
    When [action]
    Then [outcome]
```

### Godog Implementation Template
```go
func InitializeScenario(sc *godog.ScenarioContext) {
    sc.Step(`^I run "([^"]*)"$`, iRun)
    sc.Step(`^the exit code should be (\d+)$`, theExitCodeShouldBe)
}

func iRun(cmd string) error {
    // Execution logic
    return nil
}
```

## Guidelines
- **Business Value**: Focus on the "what" and "why" from the user's perspective, not the "how".
- **Reusability**: Write generic steps (e.g., "I run '...'") that can be reused across multiple scenarios.
- **Language**: Use clear, concise Gherkin that a non-technical stakeholder could understand.
