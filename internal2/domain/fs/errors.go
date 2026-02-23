package fs

import "errors"

var (
	// ErrPrecondition indicates a precondition (like ETag/CTag) check failed.
	// This usually means the item has been modified in the cloud since it was last fetched.
	ErrPrecondition = errors.New("precondition error")
)
