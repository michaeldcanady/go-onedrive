package profile

import (
	"errors"

	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
)

var (
	// ErrProfileNotFound is returned when a profile with the specified name does not exist.
	ErrProfileNotFound = errors.Join(coreerrors.ErrNotFound, errors.New("profile not found"))

	// ErrProfileAlreadyExists is returned when a profile with the specified name already exists.
	ErrProfileAlreadyExists = errors.Join(coreerrors.ErrConflict, errors.New("profile already exists"))

	// ErrProfilesBucketNotFound is returned when the profiles bucket is not found in the database.
	ErrProfilesBucketNotFound = errors.Join(coreerrors.ErrInternal, errors.New("profiles bucket not found"))

	// ErrConfigPathNotFound is returned when the configuration path for a profile is not found.
	ErrConfigPathNotFound = errors.New("config path not found")

	// ErrFailedToGetConfigDirectory is returned when the service fails to retrieve the configuration directory.
	ErrFailedToGetConfigDirectory = errors.Join(coreerrors.ErrInternal, errors.New("failed to get config directory"))

	// ErrFailedToOpenDatabase is returned when the BoltDB database cannot be opened.
	ErrFailedToOpenDatabase = errors.Join(coreerrors.ErrInternal, errors.New("failed to open database"))

	// ErrInitializeService is returned when the profile service fails to initialize.
	ErrInitializeService = errors.Join(coreerrors.ErrInternal, errors.New("failed to initialize profile service"))

	// ErrCannotDeleteDefaultProfile is returned when attempting to delete the default profile.
	ErrCannotDeleteDefaultProfile = errors.Join(coreerrors.ErrForbidden, errors.New("cannot delete the default profile"))
)
