package infra

import (
	commonerrors "github.com/michaeldcanady/go-onedrive/internal/core/errors"
)

var (
	// ErrNotFound indicates the item (file or folder) was not found.
	ErrNotFound = commonerrors.ErrNotFound
	// ErrNotFolder indicates the specified item is not a folder.
	ErrNotFolder = commonerrors.ErrNotFolder
	// ErrUnauthorized indicates the user is not authenticated.
	ErrUnauthorized = commonerrors.ErrUnauthorized
	// ErrForbidden indicates the user is authenticated but does not have
	// permission for the operation.
	ErrForbidden = commonerrors.ErrForbidden
	// ErrConflict indicates the operation failed because of a resource conflict,
	// such as a file with the same name.
	ErrConflict = commonerrors.ErrConflict
	// ErrInternal indicates an unexpected internal error.
	ErrInternal = commonerrors.ErrInternal
	// ErrPrecondition indicates a precondition (like ETag) check failed.
	ErrPrecondition = commonerrors.ErrPrecondition
	// ErrTransient indicates a temporary error that can be retried.
	ErrTransient      = commonerrors.ErrTransient
	ErrInvalidRequest = commonerrors.ErrInvalidRequest
)

type DomainError = commonerrors.DomainError

// mapGraphError2 is a lightweight variant of mapGraphError that returns only
// domain error *kinds* rather than wrapping the original error.
func mapGraphError2(err error) error {
	return mapGraphError(err, false)
}
