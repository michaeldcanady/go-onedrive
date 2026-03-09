package drive

import (
	"errors"

	commonerrors "github.com/michaeldcanady/go-onedrive/internal2/domain/common/errors"
)

var (
	ErrDriveNotFound = errors.New("drive not found")
	ErrInvalidID     = errors.New("invalid drive id")
	ErrInternal      = commonerrors.ErrInternal
	ErrNotFound      = commonerrors.ErrNotFound
	ErrUnauthorized  = commonerrors.ErrUnauthorized
	ErrForbidden     = commonerrors.ErrForbidden
	ErrConflict      = commonerrors.ErrConflict
	ErrPrecondition  = commonerrors.ErrPrecondition
	ErrTransient     = commonerrors.ErrTransient
)

type DomainError = commonerrors.DomainError
