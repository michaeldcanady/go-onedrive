# 4. Domain-First Access Pattern

Date: 2025-05-14

## Status

Status: Accepted

## Context

To ensure the CLI application is maintainable and that business rules are consistently applied, we need to prevent the UI layer (CLI commands) from directly manipulating low-level details like database records or external API calls.

## Decision

We enforce a **Domain-First Access Pattern**.

- CLI commands (the UI layer) MUST NOT access internal persistence directly (e.g., calling `bbolt` methods).
- CLI commands MUST NOT interact directly with external SDKs (e.g., `msgraph-sdk-go`) if a domain service exists for that functionality.
- Instead, commands must use **Domain Services** (e.g., `profile.Service`, `drive.Service`, `identity.Service`) which encapsulate business logic, state interaction, and external API calls.
- The UI layer is responsible for gathering user input, determining the scope of the operation (e.g., global vs session), and calling the appropriate service methods.

## Consequences

## Benefits
- **Clear Separation of Concerns:** The UI layer focuses on presentation and interaction, while domain services focus on logic and data.
- **Interchangeable Backends:** We can change the persistence layer or external API provider without modifying CLI command logic.
- **Consistent Business Logic:** Business rules (like validation or side effects) are centralized in services, ensuring they are applied regardless of which command triggers the action.
- **Improved Testability:** Domain services can be unit tested in isolation from the CLI framework.

## Trade-offs
- **Additional Abstraction:** Even simple "read" operations require a service method, which can feel like extra boilerplate for trivial tasks.
- **Learning Curve:** New contributors must understand the service boundary and avoid the temptation to take "shortcuts" by accessing data directly.

## Links

- [Domain-Driven Design (DDD) Fundamentals](https://martinfowler.com/tags/domain%20driven%20design.html)
