# GEMINI.md - Filesystem Domain

This domain defines the core filesystem abstraction for OneDrive.

## FS Domain Mandates

- **Service Interface:** All filesystem services must implement the `Service` interface in `internal2/domain/fs/service.go`.
- **Item & Metadata:** Define standardized `Item` and `Metadata` types in the domain layer.
- **Path Resolution:** Standardize the representation of OneDrive paths and their resolution to drive items.
- **Operations:** Define standard filesystem operations like `Get`, `List`, `Stat`, `ReadFile`, `WriteFile`, `Mkdir`, `Remove`, `Move`, and `Upload`.

## Engineering Standards

- **Error Handling:** Define filesystem-specific errors in `internal2/domain/fs/errors.go` or alongside the service interface.
- **Testing:** Unit tests should focus on the business logic and validation of filesystem operations, using mocks for repositories and drive resolvers.
- **Context-Aware:** All service methods must take a `context.Context` as the first argument.

## Implementation Guide

1. Define filesystem-related types like `Item`, `ListOptions`, `StatOptions`, and `ReadOptions` in `internal2/domain/fs/`.
2. Ensure the `Service` interface remains the source of truth for all filesystem operations.
3. Update the `Service` interface as needed for new filesystem features or operations.
4. Keep the filesystem domain free of infrastructure-specific details (e.g., Graph API, BoltDB).
5. Provide high-level technical rationale for the filesystem domain design.
