# Enhancing services with the Decorator pattern

This plan describes the implementation of the Decorator pattern for the
`fs.Service` and its related interfaces. The goal is to separate core provider
logic (e.g., talking to the OneDrive API) from cross-cutting concerns like
logging, validation, caching, and performance metrics.

Currently, each filesystem provider (like `onedrive` or `local`) must manually
implement its own logging and error mapping. By using decorators, we can
implement these features once and wrap any provider to add the desired
functionality.

## Objectives

The primary goals for this refactoring include:

*   **DRY (Don't Repeat Yourself):** Implement common features like structured
    logging or path validation in a single place.
*   **Composability:** Allow the assembly of "smart" services by stacking
    decorators (e.g., `Logging(Validation(Caching(OneDriveProvider)))`).
*   **Single Responsibility:** Keep provider implementations focused strictly
    on their specific storage backend logic.
*   **Consistency:** Ensure that all filesystem operations behave consistently
    regarding error handling and logging, regardless of the provider.

## Proposed infrastructure (`internal/fs/decorators`)

Decorators implement the same interface as the service they wrap, delegating the
actual work while adding their own logic before or after the call.

```go
type LoggingDecorator struct {
    base fs.Service
    log  logger.Logger
}

func (d *LoggingDecorator) Get(ctx context.Context, path *fs.URI) (fs.Item, error) {
    d.log.Debug("fetching item", logger.String("path", path.String()))
    item, err := d.base.Get(ctx, path)
    if err != nil {
        d.log.Error("failed to fetch item", logger.Error(err))
    }
    return item, err
}
```

## Potential decorators

Several key features can be extracted into decorators:

*   **`LoggingDecorator`**: Provides consistent structured logging for every
    filesystem operation.
*   **`ValidationDecorator`**: Ensures that paths and options are valid before
    calling the underlying provider.
*   **`CacheDecorator`**: Implements metadata caching to reduce API calls for
    repeated `Get` or `Stat` requests.
*   **`MetricsDecorator`**: Tracks operation latency and success rates for
    telemetry.

## Implementation steps

1.  **Define the base decorator:** Create a base struct in
    `internal/fs/decorators/` that provides default delegation for all
    `fs.Service` methods.
2.  **Implement specific decorators:** Create the `LoggingDecorator` and
    `ValidationDecorator` as the first priority.
3.  **Update the registry:** Modify `internal/fs/registry.go` to automatically
    wrap registered providers with the standard set of decorators.
4.  **Refactor providers:** Remove redundant logging and validation logic from
    the `onedrive` and `local` provider implementations.

## Next steps

Once the decorator pattern is established for the filesystem, it can be
extended to other core services like the `Identity` or `Profile` services to
provide consistent logging and validation across the entire application.
