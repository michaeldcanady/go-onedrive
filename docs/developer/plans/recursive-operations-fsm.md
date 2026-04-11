# Refactoring recursive operations with FSM

This plan describes how to refactor recursive file operations, such as `copy`
and `move`, into a state-machine-driven workflow. This change aims to improve
the reliability of large-scale operations and enable future support for
resumability.

Currently, recursive operations in `FileSystemManager` use a worker pool and
standard recursion. While efficient, this approach makes it difficult to track
fine-grained progress or recover gracefully from interruptions during
multi-gigabyte transfers involving thousands of files.

## Objectives

The primary goals for this refactoring include:

*   **Progress tracking:** Provide clear, real-time status updates for complex,
    nested directory structures.
*   **Error isolation:** Ensure that a failure in one file doesn't necessarily
    abort the entire recursive operation.
*   **Resumability:** Lay the foundation for persisting the "Work Queue" to
    the state store, allowing interrupted operations to resume where they left
    off.

## Proposed context

The `recursiveContext` struct maintains the state of the overall operation,
including the source and destination URIs and the queue of pending items.

```go
type recursiveContext struct {
    srcRoot      *shared.URI
    dstRoot      *shared.URI
    pendingItems []queueItem
    completed    int64
    totalItems   int64
    errors       []error
}

type queueItem struct {
    src *shared.URI
    dst *shared.URI
    type shared.ItemType
}
```

## Proposed states and transitions

The state machine manages the discovery and execution of the recursive task.

1.  **`ScanDirectoryState`**: Lists the contents of the current source
    directory and adds them to the `pendingItems` queue.
    *   Transitions to `ProcessQueueState`.
2.  **`ProcessQueueState`**: Picks the next item from the queue.
    *   If the item is a directory, transitions back to `ScanDirectoryState`.
    *   If the item is a file, transitions to `TransferFileState`.
    *   If the queue is empty, completes the machine execution.
3.  **`TransferFileState`**: Performs the actual file copy or move between
    providers.
    *   Transitions back to `ProcessQueueState` upon completion or error
        recording.

## Implementation steps

1.  **Define the queue:** Implement a robust queue structure that can handle
    large numbers of items.
2.  **Implement states:** Create the scanning and transfer states using the
    `pkg/fsm` package.
3.  **Update FileSystemManager:** Replace the recursive function calls in
    `copyRecursive` with the new FSM implementation.
4.  **Add progress reporting:** Integrate the state machine transitions with
    the UI layer to provide user feedback.
