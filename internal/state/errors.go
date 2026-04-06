package state

import (
	"errors"

	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
)

var (
	// ErrKeyNotFound is returned when a requested key does not exist in the state.
	ErrKeyNotFound = errors.Join(coreerrors.ErrNotFound, errors.New("state key not found"))

	// ErrBucketNotFound is returned when a required bucket is missing in the BoltDB.
	ErrBucketNotFound = errors.Join(coreerrors.ErrInternal, errors.New("state bucket not found"))
)
