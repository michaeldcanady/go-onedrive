package profile

import (
	"errors"
)

var (
	// ErrProfileNotFound is returned when a profile with the specified name does not exist.
	ErrProfileNotFound = errors.New("profile not found")

	// ErrProfileAlreadyExists is returned when a profile with the specified name already exists.
	ErrProfileAlreadyExists = errors.New("profile already exists")

	// ErrProfilesBucketNotFound is returned when the profiles bucket is not found in the database.
	ErrProfilesBucketNotFound = errors.New("profiles bucket not found")
)
