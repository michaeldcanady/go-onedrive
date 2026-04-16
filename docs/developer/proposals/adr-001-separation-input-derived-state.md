# ADR-001: Separation of Input and Derived State

## Status
Proposed

## Context
Currently, the `Options` struct in CLI commands (e.g., `internal/fs/ui/cli/rm/options.go`) stores both raw user input (like `Path string`) and objects derived from that input (like `URI *fs.URI`). This blurs the line between "data received from the user" and "data processed by the application," making it unclear when certain fields are safely populated.

## Decision
Keep the `Options` struct strictly for raw input fields and basic configuration. Move all derived or resolved state into the `CommandContext`.

```go
type CommandContext struct {
    Ctx     context.Context
    Options Options
    URI     *fs.URI // Derived from Options.Path during Validate
}
```

## Consequences
- **Pros:**
    - Clearer lifecycle: `Options` are static after `PreRunE`, `CommandContext` fields are populated during `Validate`.
    - Improved type safety: Handlers can rely on `CommandContext` fields being populated before `Execute`.
    - Easier testing: Mocking derived state is more straightforward without interfering with input parsing.
- **Cons:**
    - Slightly more verbose `CommandContext` struct.
- **Impact:** Requires a refactor of all current CLI handlers to move fields like `URI`, `SourceURI`, and `DestinationURI` from `Options` to `CommandContext`.

## Alternatives Considered
- **Option A:** Keep as is. Rejected because it leads to "partially initialized" structs being passed around.
- **Option B:** Method-based resolution. Rejected as it would require multiple calls to factory methods throughout the execution.
