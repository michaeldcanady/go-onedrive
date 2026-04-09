package drive

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/state"
)

// Service defines the interface for drive-related operations.
type Service interface {
	// ListDrives retrieves all OneDrive drives accessible to the user.
	ListDrives(ctx context.Context) ([]Drive, error)
	// ResolveDrive identifies a drive by its ID or name.
	ResolveDrive(ctx context.Context, driveRef string) (Drive, error)
	// ResolvePersonalDrive retrieves the user's primary personal OneDrive drive.
	ResolvePersonalDrive(ctx context.Context) (Drive, error)
	// GetActive retrieves the currently active drive.
	GetActive(ctx context.Context) (Drive, error)
	// SetActive marks a specific drive as the active one with the given scope.
	SetActive(ctx context.Context, driveID string, scope state.Scope) error
}
