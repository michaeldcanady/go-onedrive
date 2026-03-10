package drive

import (
	"errors"
)

var (
	ErrDriveNotFound = errors.New("drive not found")
	ErrInvalidID     = errors.New("invalid drive id")
)
