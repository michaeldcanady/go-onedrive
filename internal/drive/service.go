package drive

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/state"
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
	// GetActive retrieves the currently active drive.
	// identityID is optional and scopes the drive selection to a specific account.
	GetActive(ctx context.Context, identityID string) (Drive, error)
	// SetActive marks a specific drive as the active one with the given scope.
	// identityID is optional and scopes the drive selection to a specific account.
	SetActive(ctx context.Context, driveID string, identityID string, scope state.Scope) error
}
