# 3. Dependency injection container pattern

Date: 2025-05-14

## Status

Status: Accepted

## Context

As the number of domain services and their inter-dependencies grow, manual wiring in `main.go` or individual command initialization becomes complex, verbose, and difficult to maintain. Users need a consistent way to manage service lifecycles and inject dependencies

## Decision

Users implement a centralized **Dependency Injection (DI) Container** pattern, located in `internal/core/di/`

- The `Container` interface defines methods to retrieve some application services (for example, `Config()`, `Identity()`, `Profile()`)
- A concrete implementation of this container is responsible for instantiating and wiring services with their dependencies
- One place typically initializes services at application startup.

## Consequences

## Benefits
- **Decoupling:** Services don't need to know how to create their dependencies.
- **Testability:** Developers can mock or replace the container with a test-specific version. This makes it easier to inject mock services for unit and integration tests.
- **Single Point of Configuration:** One place centralizes all service wiring logic.
- **Consistency:** Confirms the application uses the same instance of a service throughout (where appropriate).

## Trade-offs
- **Service Locator Risk:** Passing the container around too often can lead to the Service Locator anti-pattern. Users aim to use it primarily at the application entry point and in CLI command setup.
- **Boilerplate:** Adding a new service requires updating the `Container` interface and implementation.

## Links

- [Dependency Injection in Go](https://go.dev/blog/wire) (Conceptual background)
