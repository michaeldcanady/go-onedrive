# Plan: Repository Pattern Transition

This plan outlines the steps to migrate `go-onedrive` to use a standardized Repository pattern as defined in ADR 0008.

## Phase 1: Storage Layer Standardization

- [ ] Enhance `internal/features/storage` to manage multiple database instances.
- [ ] Define common repository error types in `internal/core/errors` or `internal/core/persistence`.

## Phase 2: Feature Slice Refactoring

### Profile Slice
- [ ] Refactor `ProfileRepository` and `SettingsRepository` to accept a `*bolt.DB` handle rather than opening it.
- [ ] Update `DefaultService` to receive repositories via constructor injection.

### Identity Slice
- [ ] Rename `IdentityRepository` to `IdentityRepository` for consistency.
- [ ] Ensure the implementation follows the new standards (error handling, injection).

### Mount & Config Slices
- [ ] Refactor the `mountConfigAdapter` into a proper `MountRepository` implementation.
- [ ] Standardize interfaces and implementation for `internal/features/config`.

## Phase 3: DI and Integration

- [ ] Update `internal/core/di/service.go` to wire all repositories.
- [ ] Centralize DB lifecycle management in the `DefaultContainer`.

## Phase 4: Validation

- [ ] Ensure all unit tests pass with mocked repositories.
- [ ] Run functional tests to verify persistence still works correctly.
