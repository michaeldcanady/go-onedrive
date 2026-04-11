# Standardizing providers with the Template Method pattern

This plan describes how to use the Template Method pattern to standardize the
implementation of filesystem providers in `internal/fs/providers`. The goal is
to reduce code duplication and ensure consistent error handling across all
backends.

Currently, each provider (e.g., `onedrive`, `local`) manually implements
repetitive logic for error mapping, path normalization, and logging. By
providing a base implementation or a standard template for these operations, we
ensure that all providers follow the same architectural standards.

## Objectives

The primary goals for this refactoring include:

*   **Consistent Error Mapping:** Ensure that a `404 Not Found` from OneDrive
    and a "file not found" from the local disk are translated into the same
    domain-specific error type.
*   **Reduced Duplication:** Centralize common utility logic like path cleaning
    and URI expansion.
*   **Architectural Guardrails:** Provide a clear structure for how a provider
    should be implemented, making it easier to add new backends (e.g., S3,
    Google Drive).

## Proposed structure

We can provide a `BaseProvider` struct that implements common utility methods,
which specific providers can then embed.

```go
type BaseProvider struct {
    Log logger.Logger
}

func (b *BaseProvider) MapError(err error, uri *shared.URI) error {
    // Standard error mapping logic
}
```

## Implementation steps

1.  **Extract common logic:** Identify shared functionality across `onedrive`
    and `local` providers (e.g., error translation).
2.  **Create a base implementation:** Implement these shared methods in a
    helper file within `internal/fs/providers/shared/`.
3.  **Refactor existing providers:** Update `onedrive.Provider` and
    `local.Provider` to use the shared base logic.
4.  **Standardize decorators:** Ensure that the `Template Method` logic works
    seamlessly with the previously planned `Decorator` pattern.

## Next steps

After standardizing the provider logic, we can explore creating a "Mock
Provider" that uses the same templates, which would greatly simplify testing the
`FileSystemManager` and CLI commands.
