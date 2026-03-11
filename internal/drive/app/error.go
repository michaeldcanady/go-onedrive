package app

import (
	"errors"
)

// TODO: should this be domain layer?
var (
	ErrDriveNotFound = errors.New("drive not found")
	ErrInvalidID     = errors.New("invalid drive id")
)
