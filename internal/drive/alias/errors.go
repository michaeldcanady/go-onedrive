package alias

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
)

// NewAliasNotFoundError creates a new AppError for when an alias is not found.
func NewAliasNotFoundError(alias string) *errors.AppError {
	return errors.NewNotFound(
		nil,
		fmt.Sprintf("drive alias '%s' not found", alias),
		"Use 'odc drive alias list' to see available aliases.",
	).WithContext(errors.KeyName, alias)
}

// NewAliasAlreadyExistsError creates a new AppError for when an alias is already in use.
func NewAliasAlreadyExistsError(alias, driveID string) *errors.AppError {
	return errors.NewConflict(
		nil,
		fmt.Sprintf("drive alias '%s' already in use", alias),
		"Choose a different alias or remove the existing one.",
	).WithContext(errors.KeyName, alias).WithContext(errors.KeyDriveID, driveID)
}

// NewBucketNotFoundError creates a new AppError for when the aliases bucket is missing.
func NewBucketNotFoundError() *errors.AppError {
	return errors.NewInternal(
		nil,
		"drive aliases bucket not found",
		"This may indicate a corrupted database file.",
	)
}
