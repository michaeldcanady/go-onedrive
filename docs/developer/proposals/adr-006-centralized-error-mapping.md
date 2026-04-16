# ADR-006: Centralized Error Mapping

## Status
Proposed

## Context
CLI handlers currently return raw errors wrapped with context (e.g., `fmt.Errorf("failed to remove: %w", err)`). This makes it difficult to provide consistent exit codes (e.g., exit code 1 for generic error, 2 for missing file) or to provide localized/user-friendly error messages consistently across the CLI.

## Decision
Define a set of domain-specific errors in `internal/errors` and a centralized mapper that converts these errors into CLI-specific responses (message + exit code) during the command's `Finalize` or in a global error handler.

```go
// internal/errors
var ErrNotFound = errors.New("item not found")

// internal/fs/ui/cli/error_mapper.go
func MapToExitCode(err error) int {
    if errors.Is(err, errors.ErrNotFound) {
        return 2
    }
    return 1
}
```

## Consequences
- **Pros:**
    - Consistent UX: Users get the same exit codes for the same types of failures across all commands.
    - Separates business errors from presentation.
    - Simplifies handler logic: Handlers just return the domain error.
- **Cons:**
    - Requires careful maintenance of the error mapping table.
- **Impact:** Improved reliability for scripts that depend on `odc` exit codes.

## Alternatives Considered
- **Option A:** Set exit codes manually in each command. Rejected due to inconsistency risk.
- **Option B:** No specific mapping. Rejected as it limits the professional "feel" of the CLI tool.
