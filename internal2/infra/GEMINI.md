# GEMINI.md - Infrastructure Layer

This file outlines the mandates and standards for the `internal2/infra` package, which is the technical implementation layer of the `go-onedrive` project.

## Infrastructure Mandates

- **Technical Detail:** Implement concrete services for external systems (Microsoft Graph API, MSAL, BoltDB).
- **Domain Interfaces:** Provide implementation for repositories and low-level services defined in the `domain` layer.
- **Dependency Management:** Use the central `Container` in `internal2/app/di/container.go` for wiring and lifecycle management.
- **Error Mapping:** Map external errors (e.g., Graph API, BoltDB) to domain errors where possible.

## Engineering Standards

- **Error Handling:** Map external errors to domain-specific errors defined in the `domain` layer.
- **Logging:** Include relevant technical context (e.g., API status codes, DB paths) in the logs.
- **Testing:** Focus on the interaction with external systems, using mocks or test-specific backends where appropriate.

## Implementation Guide

When adding a new infrastructure implementation:
1. Provide a concrete implementation for a domain-defined repository or low-level service.
2. Ensure technical details like API responses, database operations, and caching are handled correctly.
3. Map any external errors to their domain equivalents in `internal2/domain/[package]/errors.go`.
4. Register the new infrastructure component in the `Container` in `internal2/app/di/container.go`.
5. Provide high-level technical rationale for the infrastructure choices.
