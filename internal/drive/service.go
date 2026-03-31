package drive

import "context"

// Service defines the interface for drive-related operations.
type Service interface {
	// ListDrives retrieves all OneDrive drives accessible to the user.
	ListDrives(ctx context.Context) ([]Drive, error)
	// ResolveDrive identifies a drive by its ID, name, or alias.
	ResolveDrive(ctx context.Context, driveRef string) (Drive, error)
	// ResolvePersonalDrive retrieves the user's primary personal OneDrive drive.
	ResolvePersonalDrive(ctx context.Context) (Drive, error)
}
