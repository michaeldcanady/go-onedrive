# 3. Dependency Injection Container Pattern

Date: 2025-05-14

## Status

Status: Accepted

## Context

As the number of domain services and their inter-dependencies grow, manual wiring in `main.go` or individual command initialization becomes complex, verbose, and difficult to maintain. We need a consistent way to manage service lifecycles and inject dependencies.

## Decision

We implement a centralized **Dependency Injection (DI) Container** pattern, located in `internal/core/di/`.

- The `Container` interface defines methods to retrieve various application services (e.g., `Config()`, `Identity()`, `Profile()`).
- A concrete implementation of this container is responsible for instantiating and wiring services with their dependencies.
- Services are typically lazily initialized or initialized once at application startup.

## Consequences

## Benefits
- **Decoupling:** Services don't need to know how their dependencies are created.
- **Testability:** The container can be mocked or replaced with a test-specific version, making it easier to inject mock services for unit and integration tests.
- **Single Point of Configuration:** All service wiring logic is centralized in one place.
- **Consistency:** Ensures that the same instance of a service is used throughout the application (where appropriate).

## Trade-offs
- **Service Locator Risk:** If the container is passed around too broadly, it can devolve into the Service Locator anti-pattern. We aim to use it primarily at the application entry point and in CLI command setup.
- **Boilerplate:** Adding a new service requires updating the `Container` interface and implementation.

## Links

- [Dependency Injection in Go](https://go.dev/blog/wire) (Conceptual background)
