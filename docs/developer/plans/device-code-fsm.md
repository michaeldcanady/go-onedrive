# Refactoring device-code authentication with FSM

This plan outlines the transition of the Microsoft device-code authentication
flow from a library-managed "black box" to a structured, state-machine-driven
implementation using the `pkg/fsm` package.

The current implementation of the device-code flow is handled entirely by the
`azidentity` library, which provides little visibility into the authentication
lifecycle. By refactoring this flow into a formal state machine, `odc` gains
greater control over the user experience, error recovery, and unit testing.

## Objectives

The primary goals for this refactoring include:

*   **Improved User Experience:** Enable the CLI to provide real-time feedback,
    such as progress spinners or countdown timers, while waiting for user
    interaction.
*   **Resilience and Auto-retry:** Automatically request a new device code if
    the current one expires, reducing the need for manual user intervention.
*   **Enhanced Testability:** Decouple the authentication steps to allow for
    isolated unit testing of the request, polling, and parsing logic using
    HTTP mocks.

## Proposed context

The `authContext` struct tracks the state and data required for the
authentication lifecycle.

```go
type authContext struct {
    tenantID        string
    clientID        string
    userCode        string
    verificationURI string
    expiresAt       time.Time
    interval        time.Duration
    deviceCode      string
    accessToken     shared.AccessToken
}
```

## Proposed states and transitions

The state machine manages the progression of the authentication flow through
discrete, well-defined states.

1.  **`RequestCodeState`**: Performs an HTTP POST to the Microsoft Entra ID
    `/devicecode` endpoint to initiate the flow.
    *   On success, transitions to `DisplayCodeState`.
    *   On failure, returns a terminal error.
2.  **`DisplayCodeState`**: Formats and displays the login instructions to the
    user (e.g., "Go to microsoft.com/devicelogin and enter code ABCD-1234").
    *   Transitions to `PollForTokenState`.
3.  **`PollForTokenState`**: Periodically checks the `/token` endpoint using
    the `device_code`.
    *   If the response is "authorization_pending", waits for the specified
        `interval` and repeats.
    *   If the code has "expired", transitions back to `RequestCodeState` to
        acquire a fresh code or returns a terminal error based on configuration.
    *   If "access_denied", returns a terminal error.
    *   On success, transitions to `PersistTokenState`.
4.  **`PersistTokenState`**: Serializes and saves the acquired token to the
    persistent state store using the `auth.SaveToken` method.
    *   On success, completes the machine execution.

## Implementation steps

Follow these steps to implement the refactored authentication flow:

1.  **Define the context:** Create the `authContext` struct and necessary
    support types in `internal/identity/providers/microsoft/`.
2.  **Implement states:** Create the individual state functions or types
    implementing the `fsm.State` interface.
3.  **Update the authenticator:** Modify the `microsoft.Authenticator` to use
    `fsm.NewMachine` with the new states instead of delegating to
    `azidentity.NewDeviceCodeCredential`.
4.  **Verify UI interaction:** Ensure the terminal output remains clear and
    helpful during the polling phase.
5.  **Add unit tests:** Implement comprehensive tests for each state, mocking
    the Microsoft Entra ID API responses.

## Next steps

After implementing the basic state machine, we can explore advanced features
such as cross-session authentication persistence, which would allow a
background process to check the status of a login initiated in a previous
terminal session.
