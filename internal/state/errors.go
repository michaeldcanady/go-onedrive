package state

import (
	"github.com/michaeldcanady/go-onedrive/internal/errors"
)

var (
	// ErrKeyNotFound is returned when a requested key does not exist in the state.
	ErrKeyNotFound = errors.NewAppError(errors.CodeNotFound, nil, "state key not found", "")

	// ErrBucketNotFound is returned when a required bucket is missing in the BoltDB.
	ErrBucketNotFound = errors.NewAppError(errors.CodeInternal, nil, "state bucket not found", "This may indicate a corrupted database file.")
)
