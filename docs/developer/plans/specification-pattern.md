# Composable filtering with the Specification pattern

This plan describes the implementation of a generic Specification pattern in
`pkg/spec`. The goal is to replace hardcoded filtering logic in `internal/fs`
with a declarative, composable system for selecting and matching items.

Currently, filtering in `internal/fs/filtering` relies on manual `if`
statements within loops. As more complex filtering criteria are added (e.g.,
by size, modification time, or extension), these loops become difficult to
maintain. The Specification pattern allows us to define small, reusable rules
that can be combined using logical operators.

## Objectives

The primary goals for this refactoring include:

*   **Composability:** Enable the creation of complex filters by combining
    simple rules (e.g., "is a file AND size > 1GB").
*   **Separation of concerns:** Move the logic for "what matches" into discrete,
    testable units, leaving the `Filterer` to only handle the iteration.
*   **Declarative syntax:** Make filtering logic more readable by expressing
    *what* should be matched rather than *how* to check for it.
*   **Extensibility:** Allow new filtering criteria to be added without
    modifying existing filter implementation logic.

## Proposed infrastructure (`pkg/spec`)

The infrastructure provides a generic interface and standard logical
combinators.

```go
type Specification[T any] interface {
    IsSatisfiedBy(item T) bool
}

// AndSpecification, OrSpecification, NotSpecification...
```

## Domain specifications example (`internal/fs/filtering`)

Specific file filters will implement the generic interface for `shared.Item`.

*   **`NamePrefixSpec`**: Matches items starting with a specific string (e.g.,
    hidden files).
*   **`TypeSpec`**: Matches items of a specific type (e.g., `TypeFolder`).
*   **`SizeSpec`**: Matches items based on a size range (e.g., greater than
    10MB).
*   **`ModifiedAfterSpec`**: Matches items modified after a specific date.

## Implementation steps

1.  **Create the engine:** Implement the generic `Specification` interface and
    combinators in `pkg/spec`.
2.  **Define domain specs:** Create a set of standard specifications for
    `shared.Item` in `internal/fs/filtering/`.
3.  **Refactor Filterer:** Update `OptionsFilterer` to use a composite
    specification instead of hardcoded `if` statements.
4.  **Update UI:** Map CLI flags (like `--all` or `--type`) to the
    corresponding specifications during command execution.

## Next steps

Once the core pattern is established, we can expand it to support more advanced
features like "Search Specifications" that can be translated into provider-
specific API queries (e.g., OData `$filter` expressions) for server-side
filtering.
