# GEMINI.md - Auth Application

This layer provides the concrete implementation of the `AuthService` interface defined in the domain layer.

## Auth Application Mandates

- **Orchestration:** Coordinate login, logout, and token acquisition.
- **Service Implementation:** Provide a robust, thread-safe implementation of `AuthService` in `internal2/app/auth/service2.go`.
- **Silent & Interactive Flows:** Support both silent (cached) and interactive (browser-based) token acquisition.
- **Caching & Profile:** Manage account records and tokens in the cache and handle profile-specific authentication.

## Engineering Standards

- **Error Handling:** Map MSAL and cache errors to domain-specific errors defined in the `domain` layer.
- **Logging:** Use `logger.WithContext(ctx)` to include correlation IDs in all auth-related logs.
- **Dependency Injection:** Register the `AuthService` in the `Container` in `internal2/app/di/container.go`.
- **Testing:** Unit tests should focus on the orchestration logic, mocking MSAL and cache dependencies.

## Implementation Guide

1. Implement the `AuthService` interface in `internal2/app/auth/service2.go`.
2. Wire the service's dependencies (Cache, Config, State, Logger, CredentialFactory, Account) in the `Container`.
3. Implement the silent and interactive token acquisition logic in `GetToken()` and `Login()`.
4. Ensure the implementation handles context-awareness and logging correctly.
5. Provide high-level technical rationale for the auth application choices.
