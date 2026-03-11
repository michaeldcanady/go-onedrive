package domain

import "errors"

var (
	// ErrDriveNotFound is returned when a requested drive is not found.
	ErrDriveNotFound = errors.New("drive not found")
	// ErrInvalidID is returned when a drive ID is invalid.
	ErrInvalidID = errors.New("invalid drive id")
)
