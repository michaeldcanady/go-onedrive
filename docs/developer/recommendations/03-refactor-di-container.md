# Recommendation: Refactor DI Container

## Goal
Decouple the dependency injection container's initialization from the concrete service implementations and move towards a more modular, provider-based approach.

## Current State
- `internal/di/service.go` contains a monolithic `NewDefaultContainer` function that manually instantiates every service.
- Adding a new service requires modifying this large function, increasing the risk of errors and making the wiring logic hard to follow.

## Value
- **High**: Improves maintainability of the application's wiring.
- Facilitates swapping implementations for testing or different environments.
- Makes the dependency graph clearer.

## Implementation Plan
1.  **Define Provider Functions**: Create small, focused functions for each service or group of services (e.g., `provideLogger`, `provideState`, `provideFileSystem`).
2.  **Lazy Initialization**: Consider using a pattern where services are initialized on first access or using a more robust DI library if the complexity warrants it (though manual provider functions are usually sufficient for Go).
3.  **Group Services**: Group related services (e.g., all filesystem-related services) into sub-providers.
4.  **Refactor `DefaultContainer`**: Update `DefaultContainer` to use these provider functions.

## Difficulty
- **Medium**: Requires a safe refactoring of the entire application's initialization sequence.
