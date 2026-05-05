# Error handling

This document describes the error handling patterns and practices used in `odc`

## Overview

`odc` uses a structured approach to error handling to provide meaningful
feedback to users and robust error recovery within the application. Errors
categorize into domain-level errors and handle consistently
across the codebase

## Standard error variables

The `internal/features/errors` package defines some standard error variables that
represent common failure scenarios:

- **`ErrNotFound`**: Indicates an item (file or folder) wasn't found
- **`ErrNotFolder`**: Indicates the specified item isn't a folder
- **`ErrUnauthorized`**: Indicates the user isn't authenticated
- **`ErrForbidden`**: Indicates the user authenticates but doesn't 
  have permission for the operation
- **`ErrConflict`**: Indicates a resource conflict, such as a file with 
  the same name
- **`ErrInternal`**: Indicates an unexpected internal error
- **`ErrPrecondition`**: Indicates a precondition (like ETag) check failed
- **`ErrTransient`**: Indicates a temporary error that can be retried
- **`ErrInvalidRequest`**: Indicates a malformed request

## Domain error pattern

The `DomainError` struct provides additional context to errors, including the
type of error, the original error, and relevant metadata like `DriveID` and
`Path`

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
  underlying error, and context
- **`Unwrap()`**: Returns the underlying error for use with `errors.Unwrap`
- **`Is(target error)`**: Reports whether the error or its kind matches the 
  target

## Error wrapping and checking

When an error occurs, it should be wrapped with context if possible. Use
`fmt.Errorf` with the `%w` verb for standard wrapping, or create a
`DomainError` for domain-specific context

### Checking errors

Always use `errors.Is` and `errors.As` for checking errors to confirm compatibility
with wrapped errors

```go
if errors.Is(err, errors.ErrNotFound) {
    // Handle not found
}

var domainErr *errors.DomainError
if errors.As(err, &domainErr) {
    // Access domainErr.Path or domainErr.DriveID
}
```

## Best practices

- **Wrap with context**: Always provide enough context to understand where and
  why the error occurred
- **Fail fast**: Validate inputs and preconditions early to avoid unnecessary
  processing
- **Meaningful messages**: Confirm error messages are clear and useful for
  the end user
- **Avoid silencing errors**: Never ignore an error; at a minimum, log it
- **Use standard errors**: Prefer the error variables in `internal/features/errors`
  when they apply
