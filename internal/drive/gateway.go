package drive

import "context"

// Gateway defines the interface for backend drive interactions (e.g., Microsoft Graph API).
type Gateway interface {
	// ListDrives fetches all available OneDrive drives.
	ListDrives(ctx context.Context, identityID string) ([]Drive, error)
	// GetPersonalDrive retrieves the user's default personal drive.
	GetPersonalDrive(ctx context.Context, identityID string) (Drive, error)
}
