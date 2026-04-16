# ADR-003: Generic Command Interface

## Status
Proposed

## Context
Each CLI command currently implements `Validate`, `Execute`, and `Finalize` methods with signatures tailored to its specific `Options` and `CommandContext`. While they follow a similar pattern, there is no formal interface, making it difficult to write shared utilities or middleware that can act on any command.

## Decision
Define a generic `Handler` interface that all CLI command implementations must satisfy.

```go
type Handler[T any] interface {
    Validate(ctx *CommandContext[T]) error
    Execute(ctx *CommandContext[T]) error
    Finalize(ctx *CommandContext[T]) error
}
```

## Consequences
- **Pros:**
    - Formalizes the "Command Pattern" across the codebase.
    - Enables generic middleware (logging, metrics, error handling).
    - Simplifies the `Create<Cmd>Cmd` boilerplate by allowing a generic runner.
- **Cons:**
    - Go's generics can sometimes add complexity to dependency injection or container wiring.
- **Impact:** Minor refactor of existing handlers to ensure they strictly adhere to the interface.

## Alternatives Considered
- **Option A:** Use `any`/interface{}. Rejected due to loss of type safety for the `Options` struct.
- **Option B:** No interface. Rejected as it prevents structural consistency.
