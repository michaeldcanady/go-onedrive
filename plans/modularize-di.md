# Architectural Improvement Plan: Modular Dependency Injection

## Goal
Simplify the central `internal/core/di` container by decentralizing service initialization logic to feature-specific packages.

## Plan
1. **Define Registry**: Define a standard `Provider` or `Initializer` interface for features to implement.
2. **Decentralize**: For each feature, move its specific wiring logic into a new `internal/features/{feature}/module.go` (or similar).
3. **Register**: Update the central DI container to simply iterate over a registered list of feature modules.
4. **Verify**: Ensure all services are initialized in the correct order to respect dependencies.

## Verification
- Test that all services are correctly injected into the application.
- Validate that adding a new feature requires minimal changes to the central DI container.
