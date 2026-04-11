# Refactoring path resolution with FSM

This plan outlines the refactoring of multi-provider path and alias resolution
into a state machine. The goal is to flatten complex nested resolution logic
into a clear, sequential flow.

Path resolution in `odc` can involve multiple steps: parsing a URI, checking
for aliases, resolving the underlying provider, and verifying connectivity. A
state machine makes these transitions explicit and easier to manage.

## Objectives

The primary goals for this refactoring include:

*   **Flatten nested logic:** Simplify the flow when aliases point to other
    providers or nested directories.
*   **Clearer error context:** Identify exactly which stage of resolution
    failed (e.g., alias not found vs. provider not registered).
*   **Centralized validation:** Ensure each resolution step is validated
    before proceeding to the next.

## Proposed context

The `resolutionContext` struct tracks the state of the URI as it is resolved
and validated across different providers.

```go
type resolutionContext struct {
    inputPath    string
    currentURI   *shared.URI
    resolvedSvc  shared.Service
    finalURI     *shared.URI
    redirects    int
}
```

## Proposed states and transitions

The state machine manages the path resolution lifecycle.

1.  **`ParseURIState`**: Breaks the input path into provider, drive, and item
    segments.
    *   On success, transitions to `ResolveAliasState`.
    *   On failure, returns a terminal error.
2.  **`ResolveAliasState`**: Checks if the URI matches a registered alias.
    *   If an alias is found, updates `currentURI` and transitions back to
        `ParseURIState` (supporting alias redirection).
    *   If no alias is found, transitions to `LoadProviderState`.
3.  **`LoadProviderState`**: Retrieves the provider implementation from the
    registry.
    *   On success, transitions to `ValidateAccessState`.
    *   On failure, returns a terminal error.
4.  **`ValidateAccessState`**: Checks if the target provider is reachable and
    authenticated.
    *   Completes the machine execution on success.

## Implementation steps

1.  **Define the context:** Create the `resolutionContext` struct and ensure
    it can track multiple redirections to prevent infinite loops.
2.  **Implement states:** Create the resolution states using the `pkg/fsm`
    package.
3.  **Update the registry:** Refactor the `Resolve` method in the registry
    to use the new FSM implementation.
4.  **Add loop detection:** Ensure that aliases cannot point to each other in
    a circular manner.
