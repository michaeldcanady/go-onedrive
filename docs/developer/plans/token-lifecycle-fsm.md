# Refactoring token lifecycle with FSM

This plan describes how to manage the lifecycle of authentication tokens using
a state machine. This change will improve token management for long-running
operations and provide a clearer path for proactive token refreshing.

A token lifecycle can be complex, involving validation, expiration checks, and
silent refreshing. Using a state machine makes these transitions explicit and
enables `odc` to handle token-related events more reliably.

## Objectives

The primary goals for this refactoring include:

*   **Manage token states:** Explicitly track tokens as valid, expiring, or
    expired.
*   **Proactive refreshing:** Initiate a background refresh when a token is
    near expiration to ensure uninterrupted operations.
*   **Centralized error handling:** Handle common authentication errors (e.g.,
    invalid grant, interaction required) in a consistent manner.

## Proposed context

The `tokenContext` struct tracks the status of an authentication token and
manages the logic for refreshing it.

```go
type tokenContext struct {
    provider     string
    accessToken  string
    refreshToken string
    expiresAt    time.Time
    isRefreshing bool
    lastError    error
}
```

## Proposed states and transitions

The state machine manages the progression of the token lifecycle.

1.  **`TokenValidState`**: Periodically checks the current token for
    expiration.
    *   If the token is nearing expiration, transitions to `TokenExpiringState`.
    *   If the token is expired, transitions to `TokenExpiredState`.
2.  **`TokenExpiringState`**: Initiates a background refresh using the
    available refresh token.
    *   On success, updates the token data and transitions to `TokenValidState`.
    *   On failure, transitions to `TokenExpiredState`.
3.  **`TokenExpiredState`**: Checks for a valid refresh token.
    *   If a refresh token is present, transitions to `RefreshState`.
    *   If no refresh token is available, transitions to `RequireAuthState`.
4.  **`RefreshState`**: Attempts to refresh the access token using the
    refresh token.
    *   On success, updates the token data and transitions to `TokenValidState`.
    *   On failure, transitions to `RequireAuthState`.
5.  **`RequireAuthState`**: Signals that manual user authentication is
    required.
    *   Completes the machine execution or enters an interactive login flow.

## Implementation steps

1.  **Define the context:** Create the `tokenContext` struct to store token-
    specific information.
2.  **Implement states:** Create the token lifecycle states using the
    `pkg/fsm` package.
3.  **Update the authenticator:** Modify the `identity` providers to use the
    new FSM implementation for token management.
4.  **Integrate with long-running tasks:** Ensure that tasks like recursive
    copying or potential background syncs are aware of the token status.
