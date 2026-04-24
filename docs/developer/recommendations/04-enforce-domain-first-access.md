# Recommendation: Enforce Domain-First Access

## Goal
Ensure all UI components (commands) interact only with domain services and avoid direct access to the application's internal state.

## Current State
- `internal/root/root.go` and some other CLI components directly use the `state.Service` for operations like setting the active profile.
- This violates the architectural principle defined in `GEMINI.md` ("Domain-First Access Pattern").

## Value
- **Medium**: Improves architectural consistency and makes the UI layer decoupled from the underlying state implementation.
- Ensures that business logic is always encapsulated within domain services.

## Implementation Plan
1.  **Audit State Access**: Use `grep` or similar tools to find all direct uses of `container.State()` in `internal/root/` and `internal/**/ui/cli/`.
2.  **Expose Domain Methods**: Ensure that `profile.Service` and other domain services have methods for all necessary operations (e.g., `SetActive(ctx, name, scope)`).
3.  **Refactor Commands**: Update all commands and the root command to use these domain service methods.

## Difficulty
- **Low**: This is a straightforward refactoring task that involves replacing direct state calls with domain service calls.
