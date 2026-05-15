package vfs

import (
	"errors"

	coreerrors "github.com/michaeldcanady/go-onedrive/internal/core/errors"
)

var (
	// ErrNotFound is returned when the requested file or directory does not exist.
	ErrNotFound = coreerrors.ErrNotFound

	// ErrAlreadyExists is returned when attempting to create a file or directory that already exists.
	ErrAlreadyExists = coreerrors.ErrAlreadyExists

	// ErrPermissionDenied is returned when the storage backend rejects an operation due to insufficient permissions.
	ErrPermissionDenied = coreerrors.ErrPermissionDenied

	// ErrNotADirectory is returned when a directory operation is attempted on a file.
	ErrNotADirectory = errors.New("not a directory")

	// ErrIsADirectory is returned when a file operation is attempted on a directory.
	ErrIsADirectory = errors.New("is a directory")

	// ErrNotEmpty is returned when attempting to remove a non-empty directory.
	ErrNotEmpty = coreerrors.ErrNotEmpty

	// ErrInvalidPath is returned when the provided path is syntactically invalid or outside the VFS root.
	ErrInvalidPath = coreerrors.ErrInvalidPath

	// ErrInternal is returned when an unexpected error occurs within the VFS or a storage plugin.
	ErrInternal = coreerrors.ErrInternal

	// ErrUnavailable is returned when the underlying storage backend or plugin is unreachable.
	ErrUnavailable = coreerrors.ErrUnavailable
)
