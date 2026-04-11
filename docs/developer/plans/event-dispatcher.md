# Decoupling UI with an Event Dispatcher

This plan describes the implementation of a generic Event Dispatcher (Observer
Pattern) in `pkg/events`. The goal is to allow core domain services (like
file systems and authenticators) to broadcast progress and status updates
without knowing about the UI implementation (CLI, progress bars, etc.).

Currently, providing feedback for long-running operations like large file
uploads or recursive copies often involves passing writers or loggers deep into
the business logic. An event-driven approach allows the UI to simply subscribe
to the events it cares about.

## Objectives

The primary goals for this refactoring include:

*   **Clean separation of concerns:** Keep domain logic "pure" and focused on
    the task, while the UI layer handles presentation.
*   **Rich UI feedback:** Enable the CLI to easily implement progress bars,
    spinners, and multi-line status updates by listening to domain events.
*   **Improved testability:** Allow unit tests to verify that specific events
    were emitted with the correct data, without needing to capture stdout.
*   **Extensibility:** Easily add new observers (e.g., a file-based audit log)
    without modifying the core logic.

## Proposed infrastructure (`pkg/events`)

The infrastructure provides a thread-safe, generic way to dispatch and
subscribe to events.

```go
type Event interface {
    Name() string
}

type Handler func(Event)

type Dispatcher struct {
    handlers map[string][]Handler
}
```

## Domain events example (`internal/fs`)

Specific features will define their own events to represent significant
lifecycle moments.

*   **`FileTransferStarted`**: Emitted when a copy or upload begins.
*   **`FileTransferProgress`**: Emitted periodically with bytes transferred
    and total size.
*   **`FileTransferCompleted`**: Emitted when the operation finishes.
*   **`TaskStatusUpdated`**: Emitted when a long-running task changes its
    internal state.

## Implementation steps

1.  **Create the engine:** Implement a thread-safe `Dispatcher` in `pkg/events`.
2.  **Define events:** Create a set of standard events for file operations in
    `internal/fs/`.
3.  **Instrument the code:** Update the `writeLargeFile` FSM and recursive
    operations to emit events via the dispatcher.
4.  **Subscribe in the CLI:** Update the CLI command handlers to subscribe to
    these events and update the terminal display (e.g., using a progress bar
    library).

## Next steps

Once the core event dispatcher is in place, we can expand it to support
asynchronous dispatching for non-blocking UI updates and event filtering for
high-frequency progress events.
