# Recommendation: Improve Test Coverage

## Goal
Significantly increase unit and integration test coverage across core packages to ensure stability and facilitate safe refactoring.

## Current State
- Test coverage is currently concentrated in `pkg/spec`, `pkg/ignore`, and `internal/fs/filtering`.
- Core packages like `internal/fs`, `internal/drive`, `internal/profile`, and `internal/di` have little to no unit tests.
- High-level orchestrators like `FileSystemManager` lack tests for cross-provider operations.

## Value
- **Very High**: Critical for ensuring that new changes don't break existing functionality (regression testing).
- Essential for the "Validate" phase of the development lifecycle.
- Enables confident refactoring of complex components like the `OneDrive Provider`.

## Implementation Plan
1.  **Identify Critical Paths**: Prioritize testing for `internal/fs` (Manager and Providers) and `internal/profile`.
2.  **Mocking Dependencies**: Use `stretchr/testify/mock` or similar to create mocks for interfaces (e.g., `drive.Gateway`).
3.  **Table-Driven Tests**: Adopt Go's idiomatic table-driven testing pattern for better coverage of edge cases.
4.  **Integration Tests**: Add integration tests that use a real (but temporary) BoltDB for `internal/profile`.
5.  **Cross-Provider Tests**: Implement tests for `FileSystemManager` that use mock providers to simulate cross-provider copy and move operations.

## Difficulty
- **Medium**: Requires writing significant amounts of test code and potentially refactoring some components to make them more testable (e.g., injecting interfaces instead of concrete types).
