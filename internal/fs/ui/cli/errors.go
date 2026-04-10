package cli

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/fs/providers/local"
	"github.com/michaeldcanady/go-onedrive/internal/fs/providers/onedrive"
	"github.com/michaeldcanady/go-onedrive/internal/identity/providers/microsoft"
)

// WrapError converts provider-specific errors into actionable AppErrors for the CLI.
// it uses a discovery process to identify both the high-level operation and the root cause.
func WrapError(err error, path string) error {
	if err == nil {
		return nil
	}

	// 1. If it's already an AppError, we trust its context.
	var appErr *errors.AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	// 2. Discover the Operation Context
	opMsg := "filesystem operation failed"
	var odRead *onedrive.ReadError
	var odWrite *onedrive.WriteError
	var locRead *local.ReadError
	var locWrite *local.WriteError

	if errors.As(err, &odRead) {
		opMsg = fmt.Sprintf("failed to read from OneDrive path '%s'", path)
	} else if errors.As(err, &odWrite) {
		opMsg = fmt.Sprintf("failed to write to OneDrive path '%s'", path)
	} else if errors.As(err, &locRead) {
		opMsg = fmt.Sprintf("failed to read from local path '%s'", path)
	} else if errors.As(err, &locWrite) {
		opMsg = fmt.Sprintf("failed to write to local path '%s'", path)
	}

	// 3. Discover the Causal Reason
	reason := "an unexpected error occurred"
	hint := "Check your connection and try again."
	code := errors.CodeInternal

	// Identity / Auth errors
	if errors.Is(err, microsoft.ErrNotAuthenticated) || hasType[*errors.UnauthorizedError](err) {
		reason = "authentication is missing or invalid"
		hint = "Run 'odc auth login' to re-authenticate."
		code = errors.CodeUnauthorized
	} else if hasType[*errors.NotFoundError](err) {
		reason = "resource not found"
		hint = "Verify the path is correct and exists."
		code = errors.CodeNotFound
	} else if hasType[*errors.ForbiddenError](err) {
		reason = "access was denied"
		hint = "Check your permissions for this resource."
		code = errors.CodeForbidden
	} else if hasType[*errors.ConflictError](err) {
		reason = "a resource conflict occurred"
		hint = "Check if the item already exists or is locked."
		code = errors.CodeConflict
	} else if hasType[*errors.TransientError](err) {
		reason = "service is temporarily unavailable"
		hint = "Please wait a moment and try again."
		code = errors.CodeTransient
	} else if hasType[*onedrive.InsufficientStorageError](err) {
		reason = "insufficient storage space"
		hint = "Free up some space in your OneDrive account."
		code = errors.CodeWriteError
	} else if hasType[*onedrive.LockedError](err) {
		reason = "the resource is locked"
		hint = "Wait for the lock to be released or check if another process is using it."
		code = errors.CodeConflict
	}

	// 4. Construct the final message
	safeMsg := fmt.Sprintf("%s: %s", opMsg, reason)
	return errors.NewAppError(code, err, safeMsg, hint)
}

// hasType is a helper to check if an error has a specific type in its chain.
func hasType[T any](err error) bool {
	var target T
	return errors.As(err, &target)
}
