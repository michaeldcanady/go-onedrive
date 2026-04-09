package errors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Re-export standard errors functions for convenience.
var (
	Is = errors.Is
	As = errors.As
)

// Common context keys for AppError.
const (
	KeyDriveID = "drive_id"
	KeyPath    = "path"
	KeyName    = "name"
)

// AppError represents a secure, application-level error.
// It separates internal details (unsafe) from user-facing information (safe).
type AppError struct {
	// Code is the machine-readable error category.
	Code ErrorCode
	// Err is the raw underlying error (internal diagnostic data).
	Err error
	// SafeMsg is the sanitized message for the user.
	SafeMsg string
	// Hint is an actionable suggestion for the user.
	Hint string
	// Context contains additional metadata for logging.
	Context map[string]interface{}
}

// NewAppError creates a new AppError.
func NewAppError(code ErrorCode, err error, safeMsg string, hint string) *AppError {
	return &AppError{
		Code:    code,
		Err:     err,
		SafeMsg: safeMsg,
		Hint:    hint,
		Context: make(map[string]interface{}),
	}
}

// WithContext adds key-value metadata to the error context.
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithHint adds an actionable suggestion to the error.
func (e *AppError) WithHint(hint string) *AppError {
	e.Hint = hint
	return e
}

// Error returns the user-facing safe message.
// This implements the error interface while preventing accidental data leaks.
func (e *AppError) Error() string {
	return e.SafeMsg
}

// Unwrap returns the underlying error for internal inspection.
func (e *AppError) Unwrap() error {
	return e.Err
}

// Is reports whether the error matches the target.
// It checks if the target is an ErrorCode, another *AppError, or if the internal error matches.
func (e *AppError) Is(target error) bool {
	if code, ok := target.(ErrorCode); ok {
		return e.Code == code
	}
	var ae *AppError
	if errors.As(target, &ae) {
		return e.Code == ae.Code
	}
	return errors.Is(e.Err, target)
}

// LogFields returns structured fields for logging the internal error details safely.
func LogFields(err error) []logger.Field {
	var e *AppError
	if errors.As(err, &e) {
		fields := []logger.Field{
			logger.String("error_code", e.Code.String()),
			logger.String("safe_message", e.SafeMsg),
		}
		if e.Err != nil {
			fields = append(fields, logger.Error(e.Err))
		}
		for k, v := range e.Context {
			fields = append(fields, logger.Any(k, v))
		}
		return fields
	}
	return []logger.Field{logger.Error(err)}
}

// Format generates a rich, sanitized error message for CLI output.
func Format(err error) string {
	var e *AppError
	if errors.As(err, &e) {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Error: %s\n", e.SafeMsg))
		if e.Hint != "" {
			sb.WriteString(fmt.Sprintf("Hint:  %s\n", e.Hint))
		}
		return sb.String()
	}
	return fmt.Sprintf("Error: %v\n", err)
}
