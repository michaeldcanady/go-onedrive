# GEMINI.md - Auth Infrastructure

This layer provides the technical implementation for authentication, using Microsoft MSAL (Microsoft Authentication Library).

## Auth Infrastructure Mandates

- **MSAL SDK:** Use the Microsoft MSAL for Go SDK to handle OAuth2 and MSAL-specific authentication.
- **Credential Provider:** Implement the `CredentialProvider` interface to bridge MSAL and the application layer.
- **Token Provider:** Use a custom `TokenProvider` to manage cached MSAL tokens.
- **Credential Factory:** Provide a factory to create different types of MSAL credentials (e.g., interactive browser, shared token cache).

## Engineering Standards

- **Error Mapping:** Map MSAL SDK errors to domain-specific errors defined in the `domain` layer.
- **Logging:** Include MSAL-specific technical context (e.g., tenant ID, client ID, scopes) in the logs.
- **Testing:** Focus on the interaction with MSAL and the Graph API, using mocks or test-specific backends where appropriate.

## Implementation Guide

1. Provide a concrete implementation for the `CredentialProvider` in `internal2/infra/auth/msal/`.
2. Implement the `TokenProvider` and `CredentialFactory` as needed for different authentication methods.
3. Map any MSAL SDK errors to their domain equivalents in `internal2/domain/auth/errors.go`.
4. Register the new auth infrastructure components in the `Container` in `internal2/app/di/container.go`.
5. Provide high-level technical rationale for the auth infrastructure choices.
