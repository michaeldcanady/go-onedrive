package state

import "errors"

var (
	// ErrKeyNotFound is returned when a requested key does not exist in the state.
	ErrKeyNotFound = errors.New("state key not found")

	// ErrBucketNotFound is returned when a required bucket is missing in the BoltDB.
	ErrBucketNotFound = errors.New("state bucket not found")
)
