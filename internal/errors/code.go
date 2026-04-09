package errors

// ErrorCode defines a set of application-level error categories.
type ErrorCode int

const (
	// CodeUnknown represents an unspecified error.
	CodeUnknown ErrorCode = iota
	// CodeInternal represents an unexpected internal failure.
	CodeInternal
	// CodeNotFound represents a missing resource.
	CodeNotFound
	// CodeUnauthorized represents an authentication failure.
	CodeUnauthorized
	// CodeForbidden represents a permission failure.
	CodeForbidden
	// CodeInvalidInput represents malformed or invalid user input.
	CodeInvalidInput
	// CodeReadError represents a failure while reading data (e.g., from disk).
	CodeReadError
	// CodeWriteError represents a failure while writing data (e.g., to disk).
	CodeWriteError
	// CodeInvalidConfig represents a configuration error.
	CodeInvalidConfig
	// CodeConflict represents a resource conflict.
	CodeConflict
	// CodePrecondition represents a failed precondition.
	CodePrecondition
	// CodeTimeout represents a timeout failure.
	CodeTimeout
	// CodeCanceled represents an operation that was canceled.
	CodeCanceled
	// CodeTransient represents a temporary failure that can be retried.
	CodeTransient
)

// Error implements the error interface for ErrorCode.
func (c ErrorCode) Error() string {
	return c.String()
}

// String returns the string representation of the ErrorCode.
func (c ErrorCode) String() string {
	switch c {
	case CodeInternal:
		return "internal_error"
	case CodeNotFound:
		return "not_found"
	case CodeUnauthorized:
		return "unauthorized"
	case CodeForbidden:
		return "forbidden"
	case CodeInvalidInput:
		return "invalid_input"
	case CodeReadError:
		return "read_error"
	case CodeWriteError:
		return "write_error"
	case CodeInvalidConfig:
		return "invalid_config"
	case CodeConflict:
		return "conflict"
	case CodePrecondition:
		return "precondition_failed"
	case CodeTimeout:
		return "timeout"
	case CodeCanceled:
		return "canceled"
	case CodeTransient:
		return "transient_error"
	default:
		return "unknown_error"
	}
}
