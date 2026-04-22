package drive

import (
	"context"
)

// Service defines the interface for drive-related operations.
type Service interface {
	// ListDrives retrieves all OneDrive drives accessible to the user.
	// identityID is optional and scopes the request to a specific account.
	ListDrives(ctx context.Context, identityID string) ([]Drive, error)
	// ResolveDrive identifies a drive by its ID, name, or alias.
	// identityID is optional and scopes the request to a specific account.
	ResolveDrive(ctx context.Context, driveRef string, identityID string) (Drive, error)
	// ResolvePersonalDrive retrieves the user's primary personal OneDrive drive.
	// identityID is optional and scopes the request to a specific account.
	ResolvePersonalDrive(ctx context.Context, identityID string) (Drive, error)
}
