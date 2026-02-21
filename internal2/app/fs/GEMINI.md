# GEMINI.md - Filesystem Application

This layer provides the concrete implementation of the `Service` interface defined in the domain layer.

## FS Application Mandates

- **Orchestration:** Coordinate filesystem operations between the domain and infrastructure layers.
- **Service Implementation:** Provide a robust, thread-safe implementation of `Service` in `internal2/app/fs/service.go`.
- **Path Resolution:** Translate human-readable paths to OneDrive item and drive IDs using the `DriveResolver`.
- **Metadata & Contents:** Manage item metadata and file contents through the `MetadataRepository` and `ContentsRepository`.
- **Caching:** Implement caching at the repository layer, not the application layer, for better control over metadata and file contents.

## Engineering Standards

- **Error Handling:** Map Graph API and cache errors to domain-specific errors defined in the `domain` layer.
- **Logging:** Use `logger.WithContext(ctx)` to include correlation IDs in all filesystem-related logs.
- **Dependency Injection:** Register the `Service` in the `Container` in `internal2/app/di/container.go`.
- **Testing:** Unit tests should focus on the orchestration logic, mocking repositories and drive resolvers.

## Implementation Guide

1. Implement the `Service` interface in `internal2/app/fs/service.go`.
2. Wire the service's dependencies (MetadataRepository, ContentsRepository, DriveResolver, Logger) in the `Container`.
3. Implement the filesystem operations like `Get`, `List`, `ReadFile`, `WriteFile`, and `Upload`.
4. Ensure the implementation handles context-awareness and logging correctly.
5. Provide high-level technical rationale for the filesystem application choices.
