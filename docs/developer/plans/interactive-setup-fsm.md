# Refactoring interactive setup with FSM

This plan describes the refactoring of interactive configuration and profile
setup into a state-machine-driven workflow. This change will decouple user
interface interactions from the underlying configuration logic.

Commands such as `profile create` and `auth login` often involve multiple
steps where the user must provide information based on previous responses. A
state machine makes these transitions easier to manage and test without
manually simulating user input.

## Objectives

The primary goals for this refactoring include:

*   **Decouple UI from logic:** Separate the orchestration of the setup flow
    from the actual data entry and validation.
*   **Support complex flows:** Easily manage multi-step processes where the
    next prompt depends on the user's previous selection.
*   **Testable interactions:** Enable mocking of the user response at each
    state, allowing for automated testing of entire configuration wizards.

## Proposed context

The `setupContext` struct tracks the user's progress through the interactive
setup process.

```go
type setupContext struct {
    profileName string
    provider    string
    authMethod  string
    configData  map[string]string
    isValidated bool
}
```

## Proposed states and transitions

The state machine manages the progression of the interactive setup wizard.

1.  **`PromptProfileNameState`**: Requests the user to enter a name for the new
    profile.
    *   On success, transitions to `SelectProviderState`.
2.  **`SelectProviderState`**: Displays a list of available providers (e.g.,
    Microsoft, Local) and asks the user to choose one.
    *   Transitions to `SelectAuthMethodState`.
3.  **`SelectAuthMethodState`**: Prompts the user to select an authentication
    method based on the chosen provider (e.g., Interactive, Client Secret).
    *   Transitions to `CollectAuthDataState`.
4.  **`CollectAuthDataState`**: Requests the specific credentials or
    information required for the chosen authentication method.
    *   Transitions to `ValidateSetupState`.
5.  **`ValidateSetupState`**: Performs a test authentication or configuration
    check to ensure everything is correct.
    *   Completes the machine execution on success.

## Implementation steps

1.  **Define the context:** Create the `setupContext` struct and ensure it
    can store information collected during the wizard.
2.  **Implement states:** Create the interactive prompt states using the
    `pkg/fsm` package and a suitable terminal prompting library.
3.  **Update CLI commands:** Refactor the `profile create` and `auth login`
    handlers to use the new FSM implementation.
4.  **Verify UI interaction:** Ensure the terminal output remains intuitive
    and helpful throughout the interactive process.
