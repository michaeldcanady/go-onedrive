# Error Handling

This document describes the error handling patterns and practices used in `odc`.

## Overview

`odc` uses a structured approach to error handling to provide meaningful
feedback to users and robust error recovery within the application. Errors
are categorized into domain-level errors and are handled consistently
across the codebase.

## Standard Error Variables

The `internal/errors` package defines several standard error variables that
represent common failure scenarios:

- **`ErrNotFound`**: Indicates an item (file or folder) was not found.
- **`ErrNotFolder`**: Indicates the specified item is not a folder.
- **`ErrUnauthorized`**: Indicates the user is not authenticated.
- **`ErrForbidden`**: Indicates the user is authenticated but does not 
  have permission for the operation.
- **`ErrConflict`**: Indicates a resource conflict, such as a file with 
  the same name.
- **`ErrInternal`**: Indicates an unexpected internal error.
- **`ErrPrecondition`**: Indicates a precondition (like ETag) check failed.
- **`ErrTransient`**: Indicates a temporary error that can be retried.
- **`ErrInvalidRequest`**: Indicates a malformed request.

## Domain Error Pattern

The `DomainError` struct provides additional context to errors, including the
type of error, the original error, and relevant metadata like `DriveID` and
`Path`.

```go
type DomainError struct {
	Kind    error
	Err     error
	DriveID string
	Path    string
}
```

### Methods

- **`Error()`**: Returns a formatted error message including the kind, 
  underlying error, and context.
- **`Unwrap()`**: Returns the underlying error for use with `errors.Unwrap`.
- **`Is(target error)`**: Reports whether the error or its kind matches the 
  target.

## Error Wrapping and Checking

When an error occurs, it should be wrapped with context if possible. Use
`fmt.Errorf` with the `%w` verb for standard wrapping, or create a
`DomainError` for domain-specific context.

### Checking Errors

Always use `errors.Is` and `errors.As` for checking errors to ensure compatibility
with wrapped errors.

```go
if errors.Is(err, errors.ErrNotFound) {
    // Handle not found
}

var domainErr *errors.DomainError
if errors.As(err, &domainErr) {
    // Access domainErr.Path or domainErr.DriveID
}
```

## Best Practices

- **Wrap with context**: Always provide enough context to understand where and
  why the error occurred.
- **Fail fast**: Validate inputs and preconditions early to avoid unnecessary
  processing.
- **Meaningful messages**: Ensure error messages are clear and actionable for
  the end user.
- **Avoid silencing errors**: Never ignore an error; at a minimum, log it.
- **Use standard errors**: Prefer the error variables in `internal/errors`
  when they apply.
