# GEMINI.md - Application Layer

This file outlines the mandates and standards for the `internal2/app` package, which is the orchestration layer of the `go-onedrive` project.

## Application Layer Mandates

- **Orchestration:** Coordinate logic between the domain and infrastructure layers.
- **Service Implementation:** Provide concrete implementations for domain interfaces.
- **Dependency Management:** Use the central `Container` for all wiring and lifecycle management in `internal2/app/di/container.go`.
- **Concurrency:** Implement thread-safe operations where necessary, especially for caching and state management.

## Engineering Standards

- **Error Handling:** Wrap errors from the infrastructure layer with `fmt.Errorf("...: %w", err)`.
- **Logging:** Use `logger.WithContext(ctx)` to ensure correlation IDs are propagated.
- **Dependency Injection:** Use the container's lazy initialization pattern with `sync.Once`.
- **Testing:** Focus on use-case validation, mocking infrastructure dependencies (e.g., API, DB, Cache).

## Implementation Guide

When adding a new service implementation:
1. Implement the domain interface in `internal2/app/[package]/service_impl.go`.
2. Add the service and its `sync.Once` to the `Container` in `internal2/app/di/container.go`.
3. Wire the service's dependencies by calling other container methods.
4. Ensure the implementation handles context-awareness and logging correctly.
5. Provide high-level technical rationale for the implementation choices.
