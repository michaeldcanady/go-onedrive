package mount

import "context"

// Mount represents a logical attachment point that maps a VFS path segment
// to a specific storage backend and an authenticated identity.
type Mount struct {
	Path             string            `json:"path"`
	Type             string            `json:"type"` // e.g., "local", "onedrive", "googledrive"
	IdentityID       string            `json:"identity_id"`
	IdentityProvider string            `json:"identity_provider"` // e.g., "azure", "google"
	Options          map[string]string `json:"options"`
}

// Service coordinates the lifecycle of [Mount] points within the application.
type Service interface {
	// Add registers a new mount point and persists it to the repository.
	Add(ctx context.Context, m *Mount) error

	// List returns all currently registered mount points.
	List(ctx context.Context) ([]*Mount, error)

	// Remove unregisters the mount point at the specified path.
	Remove(ctx context.Context, path string) error

	// Get retrieves the mount point configuration for the specified path.
	Get(ctx context.Context, path string) (*Mount, error)
}

// Repository handles the low-level persistence of [Mount] configurations.
type Repository interface {
	// Save persists the mount configuration to the underlying store.
	Save(m *Mount) error

	// List retrieves all stored mount configurations.
	List() ([]*Mount, error)

	// Delete removes a mount configuration by its logical path.
	Delete(path string) error

	// Get retrieves a single mount configuration by its logical path.
	Get(path string) (*Mount, error)
}
