# ADR-002: Standardized Command Result Object

## Status
Proposed

## Context
CLI commands currently communicate success by printing directly to `Stdout` in the `Finalize` method. This makes it difficult to reuse command logic programmatically, captures output for unit tests via buffer interception (which is brittle), and doesn't provide a structured way to handle different output formats (JSON, YAML) consistently.

## Decision
Introduce a `Result` object that `Execute` returns (or populates in `CommandContext`). `Finalize` should then be responsible for rendering this `Result` object.

```go
type Result struct {
    Message string      // Human-readable summary
    Data    any         // Structured data for JSON/YAML output
    Status  int         // Exit code or specific status identifier
}
```

## Consequences
- **Pros:**
    - Programmatic access: Commands can be called by other internal services and return structured data.
    - Consistent formatting: Rendering logic can be centralized or shared across commands.
    - Improved testability: Assertions can be made against the `Result` object instead of scraping stdout.
- **Cons:**
    - Requires defining `Result` structures for every command.
- **Impact:** Changes the signature of `Execute` or adds a `Result` field to `CommandContext`.

## Alternatives Considered
- **Option A:** Continue printing in `Finalize`. Rejected as it limits the CLI's flexibility for future integrations (like an API or GUI).
- **Option B:** Return only error. Rejected because success metadata is often needed.
