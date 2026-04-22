# Architectural Improvement Plan: Consolidate Domain Errors

## Goal
Enforce consistent error handling by mandating the use of `internal/core/errors.DomainError` across all feature slices.

## Plan
1. **Audit Errors**: Identify all locations where non-`DomainError` errors are being used in feature slices.
2. **Replace**: Replace custom or primitive error types with `DomainError` (or wrap them using `errors.Wrap`).
3. **Standardize**: Update the documentation/linting rules to discourage the creation of new feature-local error types.

## Verification
- Confirm that consistent, user-friendly error messages are output by the CLI.
- Ensure that the error handling logic in the CLI layer remains intact and compatible with the consolidated `DomainError` type.
