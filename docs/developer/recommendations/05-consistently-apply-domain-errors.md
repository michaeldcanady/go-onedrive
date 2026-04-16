# Recommendation: Consistently Apply Domain Errors

## Goal
Standardize error handling across the application by ensuring all services return domain-specific errors.

## Current State
- `internal/errors` contains custom error types (`DomainError`), but their usage is inconsistent across the codebase.
- Many services return raw errors or `fmt.Errorf` without wrapping them in domain-specific kinds (e.g., `ErrNotFound`, `ErrUnauthorized`).

## Value
- **Medium**: Improves the predictability of error handling for callers.
- Enables more refined error reporting to the user (e.g., distinguishing between a file not found and a permission error).

## Implementation Plan
1.  **Define Error Kinds**: Audit the existing error kinds in `internal/errors` and add any missing ones (e.g., `ErrDriveNotFound`, `ErrInvalidProfile`).
2.  **Audit Service Returns**: Systematic check of all service methods in `internal/` to ensure errors are mapped or wrapped into `DomainError`.
3.  **Standardize Mapping**: Use common error mapping helpers for frequent external errors (e.g., from the Microsoft Graph SDK).

## Difficulty
- **Low/Medium**: Requires an exhaustive check of all service methods and consistent application of wrapping.
