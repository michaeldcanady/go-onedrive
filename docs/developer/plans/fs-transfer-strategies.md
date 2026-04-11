# Decoupling file operations with the Strategy pattern

This plan describes how to refactor the `FileSystemManager` to use the Strategy
pattern for core operations like `Copy` and `Move`. The goal is to separate the
decision-making logic (determining the type of operation) from the execution
details (performing the actual transfer).

Currently, the `FileSystemManager` contains branching logic to handle same-
provider vs. cross-provider transfers. As more providers and specialized
transfer techniques (e.g., parallel multipart uploads) are added, this logic
will become increasingly complex. The Strategy pattern allows us to encapsulate
each transfer method into its own dedicated class.

## Objectives

The primary goals for this refactoring include:

*   **Clean branching logic:** Move complex `if/else` checks for provider
    compatibility into a factory that selects the appropriate strategy.
*   **Specialized optimizations:** Easily implement provider-specific
    optimizations (e.g., server-side copy for OneDrive) without cluttering the
    general manager logic.
*   **Extensibility:** Add support for new transfer methods (e.g., streaming
    transfers, peer-to-peer) simply by adding a new strategy implementation.
*   **Improved testing:** Allow individual transfer strategies to be unit-
    tested in isolation with mock providers.

## Proposed infrastructure (`internal/fs/strategies`)

The infrastructure defines a common interface for file transfer strategies.

```go
type TransferStrategy interface {
    Execute(ctx context.Context, src, dst *shared.URI, opts shared.CopyOptions) error
}
```

## Potential strategies

Several specialized strategies can be implemented:

*   **`InternalTransferStrategy`**: Used when the source and destination are on
    the same provider, allowing for optimized server-side operations.
*   **`CrossProviderTransferStrategy`**: Handles streaming data between
    different providers (e.g., downloading from Local and uploading to OneDrive).
*   **`MultipartTransferStrategy`**: A specialized strategy for very large
    files that uses parallel chunked uploads for better performance.

## Implementation steps

1.  **Define the strategy interface:** Create the `TransferStrategy` interface
    in `internal/fs/strategies/`.
2.  **Implement core strategies:** Move the current same-provider and cross-
    provider logic from `manager.go` into dedicated strategy implementations.
3.  **Create a strategy factory:** Implement a component that inspects the
    source and destination URIs and returns the most efficient strategy.
4.  **Refactor the manager:** Update `FileSystemManager` to use the factory and
    selected strategy for all `Copy` and `Move` operations.

## Next steps

Once the strategy pattern is in place, we can easily add a "Dry Run" strategy
that logs what *would* happen without performing any actual file operations,
providing a safe way for users to preview complex recursive tasks.
