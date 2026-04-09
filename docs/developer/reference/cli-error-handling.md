# Error Handling

This document describes the error handling patterns and practices used in `odc`.

## Overview

`odc` uses a structured approach to error handling to provide meaningful, secure, and actionable feedback to users while maintaining robust internal diagnostics. The core of this system is the `AppError` type.

## AppError Design

The `AppError` type separates internal error details (unsafe for users) from sanitized, user-facing information.

```go
type AppError struct {
	Code    ErrorCode              // Machine-readable category
	Err     error                  // Raw underlying error (for logging)
	SafeMsg string                 // Sanitized message for the user
	Hint    string                 // Actionable suggestion
	Context map[string]interface{} // Metadata (DriveID, Path, etc.)
}
```

### Key Principles

1.  **Security**: Internal error messages (e.g., database connection strings, raw API responses) are kept in the `Err` field and logged, but never shown to the user.
2.  **Actionability**: Every error should ideally include a `Hint` telling the user how to fix the problem.
3.  **Categorization**: Errors are categorized using `ErrorCode` (enum) for consistent handling and reporting.

## Error Codes

The `internal/errors` package defines several error categories:

- `CodeInternal`: Unexpected internal failures.
- `CodeNotFound`: Missing resources (files, profiles, etc.).
- `CodeUnauthorized`: Authentication failures.
- `CodeForbidden`: Permission failures.
- `CodeInvalidInput`: Malformed user input.
- `CodeConflict`: Resource conflicts (e.g., file already exists).
- `CodeTransient`: Temporary failures (e.g., network timeout, service busy).
- `CodeReadError` / `CodeWriteError`: Disk I/O failures.

## Usage in Services

Services should return `AppError` when they have enough context to provide a safe message or hint.

### Creating Errors

Use the convenience functions in `internal/errors`:

```go
import "github.com/michaeldcanady/go-onedrive/internal/errors"

// Simple creation
return errors.NewNotFound(err, "profile not found", "Use 'odc profile list' to see available profiles.")

// With context
return errors.NewAppError(errors.CodeWriteError, err, "could not save config", "")
    .WithContext(errors.KeyPath, path)
```

### Standard Context Keys

- `errors.KeyPath`: The file or OneDrive path involved.
- `errors.KeyDriveID`: The OneDrive Drive ID.
- `errors.KeyName`: A resource name (e.g., profile name).

## Usage in CLI Handlers

CLI handlers follow a strict "Log and Return" pattern.

```go
func (h *Handler) Handle(ctx context.Context, opts Options) error {
    res, err := h.service.DoSomething(ctx)
    if err != nil {
        // 1. Log with full context
        h.log.Error(err.Error(), errors.LogFields(err)...)
        
        // 2. Return the error (main.go will handle formatting and printing to stderr)
        return err
    }
    
    // ... success path
    return nil
}
```

**Note:** Handlers should **never** print errors to `Stdout` or `Stderr` themselves if they are returning them. `main.go` uses `errors.Format(err)` to ensure consistent, beautiful error output for all commands.

## Checking Errors

Use `errors.Is` to check for specific error categories or underlying errors.

```go
// Check by category
if errors.Is(err, errors.CodeNotFound) { ... }

// Check for a specific AppError instance (e.g. from a service)
if errors.Is(err, profile.ErrProfileNotFound) { ... }
```

## Legacy Support

The `DomainError` type is deprecated and being phased out in favor of `AppError`. New code should always use `AppError`.
