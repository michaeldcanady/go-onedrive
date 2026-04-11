# Standardizing URI creation with the Factory pattern

This plan describes the implementation of the Factory pattern for creating and
validating `shared.URI` objects. The goal is to centralize the parsing logic
and ensure that all URIs used within `odc` are consistent and correctly
formatted.

Currently, URI parsing is handled by a single `ParseURI` function in
`internal/fs/uri.go`. As we add support for more providers and complex path
structures (e.g., drive aliases, absolute paths, relative paths), this function
will become increasingly complex. A dedicated `URIFactory` allows us to manage
these different creation rules in a structured way.

## Objectives

*   **Centralized Parsing Logic:** Move all path parsing and normalization rules
    into a single, dedicated component.
*   **Support for Multiple Formats:** Easily handle different URI formats
    (e.g., `local:/tmp`, `onedrive:me:/Documents`, `alias:path`) through
    specialized factory methods.
*   **Validation at Creation:** Ensure that every `shared.URI` object is
    syntactically valid and points to a supported provider at the moment it is
    instantiated.
*   **Decoupled URI Construction:** Separate the high-level request (e.g.,
    "parse this user input") from the low-level implementation details of URI
    structure.

## Proposed infrastructure (`internal/fs/uri_factory.go`)

The factory provides a set of methods for creating URIs from different sources.

```go
type URIFactory struct {
    registry ServiceRegistry
    aliasSvc alias.Service
}

func (f *URIFactory) FromString(input string) (*shared.URI, error) { ... }
func (f *URIFactory) FromLocalPath(path string) (*shared.URI, error) { ... }
func (f *URIFactory) FromAlias(name, subpath string) (*shared.URI, error) { ... }
```

## Implementation steps

1.  **Define the factory:** Create the `URIFactory` struct and its methods in
    `internal/fs/uri_factory.go`.
2.  **Migrate parsing logic:** Move the logic from the current `ParseURI`
    function into the factory's `FromString` method.
3.  **Integrate with DI:** Add the `URIFactory` to the `DefaultContainer` and
    ensure it has access to the provider registry and alias service.
4.  **Refactor CLI commands:** Update command handlers to use the factory for
    all path processing, ensuring consistent validation across the application.

## Next steps

Once the URI factory is established, we can expand it to support "Contextual
Resolution," where the factory can automatically resolve relative paths or
shorthand notations based on the user's current "working directory" in
OneDrive.
