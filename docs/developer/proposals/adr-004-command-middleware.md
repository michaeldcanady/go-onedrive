# ADR-004: Command Middleware

## Status
Proposed

## Context
Standard operational logic—such as logging "starting operation," timing execution for metrics, and catching panics—is currently implemented manually (and often inconsistently) within the `Execute` method of each command. This leads to boilerplate duplication and makes it hard to change cross-cutting concerns.

## Decision
Implement a "Middleware" pattern for CLI commands. A generic runner will wrap the `Execute` call with a chain of middleware functions.

```go
type Middleware func(next CommandFunc) CommandFunc

func LoggingMiddleware(next CommandFunc) CommandFunc {
    return func(ctx *CommandContext) error {
        log.Info("Starting operation")
        err := next(ctx)
        log.Info("Operation finished")
        return err
    }
}
```

## Consequences
- **Pros:**
    - DRY (Don't Repeat Yourself): Cross-cutting concerns are implemented once.
    - Consistency: Every command gets logging, metrics, and error handling for "free."
    - Clean Business Logic: Handlers focus only on the actual operation.
- **Cons:**
    - Adds a layer of abstraction that might make debugging slightly more complex.
- **Impact:** Significant reduction in boilerplate within `handler.go` files.

## Alternatives Considered
- **Option A:** Manual implementation in every handler. Rejected due to maintainability costs.
- **Option B:** Inheritance/Base struct. Rejected as Go prefers composition over inheritance.
