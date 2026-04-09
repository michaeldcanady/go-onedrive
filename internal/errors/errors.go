package errors

// NewNotFound creates a new AppError for missing resources.
func NewNotFound(err error, safeMsg, hint string) *AppError {
	return NewAppError(CodeNotFound, err, safeMsg, hint)
}

// NewInternal creates a new AppError for internal failures.
func NewInternal(err error, safeMsg, hint string) *AppError {
	return NewAppError(CodeInternal, err, safeMsg, hint)
}

// NewInvalidInput creates a new AppError for malformed user input.
func NewInvalidInput(err error, safeMsg, hint string) *AppError {
	return NewAppError(CodeInvalidInput, err, safeMsg, hint)
}

// NewUnauthorized creates a new AppError for authentication failures.
func NewUnauthorized(err error, safeMsg, hint string) *AppError {
	return NewAppError(CodeUnauthorized, err, safeMsg, hint)
}

// NewForbidden creates a new AppError for permission failures.
func NewForbidden(err error, safeMsg, hint string) *AppError {
	return NewAppError(CodeForbidden, err, safeMsg, hint)
}

// NewConflict creates a new AppError for resource conflicts.
func NewConflict(err error, safeMsg, hint string) *AppError {
	return NewAppError(CodeConflict, err, safeMsg, hint)
}
