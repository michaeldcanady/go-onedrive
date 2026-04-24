# Recommendations for Improving the go-onedrive Codebase

The following plans outline a series of recommendations for improving the maintainability, testability, and architectural consistency of the `go-onedrive` (odc) project. These plans are ordered based on their overall value and impact on the project's long-term health.

| # | Plan | Value | Difficulty | Description |
| :--- | :--- | :--- | :--- | :--- |
| 1 | [Improve Test Coverage](./01-improve-test-coverage.md) | Very High | Medium | Increase unit and integration testing across core packages. |
| 2 | [Decompose OneDrive Provider](./02-decompose-onedrive-provider.md) | High | Medium | Break down the monolithic provider into smaller components. |
| 3 | [Refactor DI Container](./03-refactor-di-container.md) | High | Medium | Modularize the dependency injection container. |
| 4 | [Enforce Domain-First Access](./04-enforce-domain-first-access.md) | Medium | Low | Ensure UI components interact only with domain services. |
| 5 | [Consistently Apply Domain Errors](./05-consistently-apply-domain-errors.md) | Medium | Low | Standardize error handling throughout the application. |
| 6 | [Refactor CLI Command Boilerplate](./06-refactor-cli-command-boilerplate.md) | Medium | Low | Reduce repetitive code in CLI command definitions. |
