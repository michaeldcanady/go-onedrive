package drive

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/state"
)

// Service defines the interface for drive-related operations.
type Service interface {
	// ListDrives retrieves all OneDrive drives accessible to the user.
	ListDrives(ctx context.Context) ([]Drive, error)
	// ResolveDrive identifies a drive by its ID, name, or alias.
	ResolveDrive(ctx context.Context, driveRef string) (Drive, error)
	// ResolvePersonalDrive retrieves the user's primary personal OneDrive drive.
	ResolvePersonalDrive(ctx context.Context) (Drive, error)
	// GetActive retrieves the currently active drive.
	GetActive(ctx context.Context) (Drive, error)
	// SetActive marks a specific drive as the active one with the given scope.
	// identityID is optional and scopes the drive selection to a specific account.
	SetActive(ctx context.Context, driveID string, identityID string, scope state.Scope) error
}
