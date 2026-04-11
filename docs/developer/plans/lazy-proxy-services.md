# Optimizing startup with Lazy-Loading Proxies

This plan describes the implementation of the Proxy pattern to enable lazy-
loading of heavy application services in `internal/di`. The goal is to improve
the startup performance of the `odc` CLI by only initializing dependencies
when they are actually required by a command.

Currently, `NewDefaultContainer` initializes all core services (including the
Microsoft Graph SDK and the BoltDB state store) every time the CLI is run. For
simple commands like `--version` or `--help`, this results in unnecessary
overhead. A lazy-loading proxy allows us to defer this initialization until the
service's methods are first invoked.

## Objectives

*   **Faster startup:** Reduce the time it takes for simple or non-cloud-related
    commands to execute by deferring expensive service setup.
*   **Resource efficiency:** Avoid opening database connections or initializing
    complex SDKs unless the specific execution path requires them.
*   **Encapsulated initialization:** Keep the service wiring in the container
    clean while managing the lifecycle of heavy dependencies behind a proxy.

## Proposed infrastructure (`internal/di/proxies`)

The infrastructure provides a wrapper that implements the service interface and
holds a reference to the real implementation, which is initialized on demand.

```go
type DriveServiceProxy struct {
    container *DefaultContainer
    real      drive.Service
    once      sync.Once
}

func (p *DriveServiceProxy) ListDrives(ctx context.Context) ([]Drive, error) {
    p.once.Do(func() {
        p.real = p.container.initRealDriveService()
    })
    return p.real.ListDrives(ctx)
}
```

## Potential proxied services

Heavy services that are good candidates for lazy-loading include:

*   **`OneDriveProvider`**: Requires the Graph SDK and authentication token
    validation.
*   **`BoltService`**: Requires opening a file lock on the state database.
*   **`GraphProvider`**: Involves complex initialization of the Microsoft
    identity and middleware stack.

## Implementation steps

1.  **Identify heavy services:** Profile the CLI startup to determine which
    components contribute most to initialization latency.
2.  **Create proxy wrappers:** Implement proxy structs for these services in a
    dedicated package or within the `di` package.
3.  **Update the container:** Modify `NewDefaultContainer` to return proxies
    instead of fully initialized service instances.
4.  **Manage shared state:** Ensure that proxies correctly handle dependencies
    on other proxied services to avoid circular or failed initializations.

## Next steps

After implementing lazy-loading proxies, we can explore "Pre-fetching" logic
that starts the initialization of cloud services in the background while the
CLI performs initial argument validation, further reducing perceived latency.
