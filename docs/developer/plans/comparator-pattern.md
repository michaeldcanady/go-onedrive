# Composable sorting with the Comparator pattern

This plan describes the implementation of a generic Comparator pattern in
`pkg/sortutil`. The goal is to replace reflection-based sorting in
`internal/fs/sorting` with a type-safe, declarative, and composable system.

Currently, sorting in `internal/fs/sorting` relies on `reflect` to access
fields by string name. This is slow, error-prone, and difficult to extend with
multi-level sorting (e.g., "sort by folders first, then by name"). The
Comparator pattern allows us to define small, reusable ordering rules that can
be chained together.

## Objectives

The primary goals for this refactoring include:

*   **Type safety:** Eliminate the use of reflection in favor of static,
    compile-time checked comparisons.
*   **Composability:** Enable multi-level sorting (e.g., `ByFolderFirst.Then(ByName)`)
    through a simple chaining API.
*   **Performance:** Improve execution speed by removing the overhead of
    reflection during the `sort.Slice` inner loop.
*   **Readability:** Make sorting logic more expressive by using named
    comparators instead of dynamic field mapping.

## Proposed infrastructure (`pkg/sortutil`)

The infrastructure provides a generic interface and standard combinators for
building complex ordering rules.

```go
type Comparator[T any] interface {
    Less(i, j T) bool
}

// Then combines two comparators to create a primary/secondary ordering.
func Then[T any](primary, secondary Comparator[T]) Comparator[T] {
    return thenComparator[T]{primary, secondary}
}
```

## Domain comparators example (`internal/fs/sorting`)

Specific item comparators will implement the generic interface for `shared.Item`.

*   **`NameComparator`**: Compares items alphabetically by name.
*   **`SizeComparator`**: Compares items by their file size.
*   **`DateComparator`**: Compares items by their last modified timestamp.
*   **`FolderFirstComparator`**: Ensures folders always appear before files.
*   **`ReverseComparator`**: Inverts the logic of any underlying comparator.

## Implementation steps

1.  **Create the engine:** Implement the generic `Comparator` interface and
    `Then` / `Reverse` combinators in `pkg/sortutil`.
2.  **Define domain comparators:** Create a set of standard comparators for
    `shared.Item` in `internal/fs/sorting/`.
3.  **Refactor Sorter:** Update `OptionsSorter` to build a composite comparator
    based on user options instead of using reflection.
4.  **Simplify comparison logic:** Remove the reflection-based `compare_values.go`
    in favor of explicit, type-safe comparison functions.

## Next steps

Once the core comparator pattern is established, it can be used to implement
more advanced sorting features, such as natural sort (e.g., "file2.txt" before
"file10.txt") or locale-aware string comparison.
