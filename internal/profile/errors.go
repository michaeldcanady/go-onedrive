package profile

import (
	"github.com/michaeldcanady/go-onedrive/internal/errors"
)

var (
	// ErrProfileNotFound is returned when a profile with the specified name does not exist.
	ErrProfileNotFound = errors.NewAppError(errors.CodeNotFound, nil, "profile not found", "Use 'odc profile list' to see available profiles.")

	// ErrProfileAlreadyExists is returned when a profile with the specified name already exists.
	ErrProfileAlreadyExists = errors.NewAppError(errors.CodeConflict, nil, "profile already exists", "Choose a different name for the profile.")

	// ErrProfilesBucketNotFound is returned when the profiles bucket is not found in the database.
	ErrProfilesBucketNotFound = errors.NewAppError(errors.CodeInternal, nil, "profiles bucket not found", "This may indicate a corrupted database file.")
)
