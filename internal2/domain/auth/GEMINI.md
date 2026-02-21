# GEMINI.md - Auth Domain

This domain handles the core authentication logic and interfaces for the `go-onedrive` project.

## Auth Domain Mandates

- **AuthService Interface:** All authentication services must implement the `AuthService` interface.
- **Token Credential:** The `AuthService` should also implement the `azcore.TokenCredential` interface for use with Azure SDKs.
- **Secure Handling:** Authentication records, tokens, and sensitive account information must be handled securely.
- **Profile-Aware:** All auth operations must be aware of the user's profile.

## Engineering Standards

- **Error Handling:** Define auth-specific errors in `internal2/domain/auth/errors.go`.
- **Testing:** Unit tests should focus on login/logout state transitions and token acquisition logic, using mocks where necessary.
- **Context-Aware:** All service methods must take a `context.Context` as the first argument.

## Implementation Guide

1. Define auth-related types like `LoginOptions`, `LoginResult`, and `AccessToken` in `internal2/domain/auth/`.
2. Ensure the `AuthService` interface remains the source of truth for all auth operations.
3. Update the `AuthService` interface as needed for new authentication methods or features.
4. Keep the auth domain free of infrastructure-specific details (e.g., MSAL, BoltDB).
5. Provide high-level technical rationale for the auth domain design.
