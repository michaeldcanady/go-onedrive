package drive

import (
	"context"
)

// Drive represents a storage container, such as a personal OneDrive,
// a SharePoint document library, or a local directory.
type Drive struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	IdentityID string `json:"identity_id"`
	Type       string `json:"type"` // e.g., "personal", "business"
}

// Service coordinates the discovery and retrieval of [Drive] metadata
// by orchestrating requests across registered storage plugins.
type Service interface {
	// List returns all drives accessible to the specified identity.
	// If identityID is empty, it attempts to list drives for all active identities.
	List(ctx context.Context, identityID string) ([]*Drive, error)

	// Get retrieves drive metadata from the local persistent cache.
	Get(ctx context.Context, driveID string) (*Drive, error)

	// FindDrive searches for a drive matching the provided query (ID or name).
	FindDrive(ctx context.Context, query string) (*Drive, error)
}

// Repository handles the persistent caching of [Drive] metadata to avoid
// redundant and expensive plugin-based discovery.
type Repository interface {
	// Save persists the drive metadata to the underlying store.
	Save(d *Drive) error

	// ListByIdentity retrieves cached drives filtered by the associated identity.
	ListByIdentity(identityID string) ([]*Drive, error)

	// ByID retrieves a single drive by its unique identifier from the cache.
	ByID(driveID string) (*Drive, error)

	// Delete removes drive metadata from the cache.
	Delete(driveID string) error
}

// IdentityService defines the subset of identity management required by the drive feature.
type IdentityService interface {
	GetIdentity(ctx context.Context, id string) (*Identity, error)
	List(ctx context.Context) ([]*Identity, error)
}

// TokenService manages access tokens for storage operations.
type TokenService interface {
	GetToken(ctx context.Context, provider, identityID string) (*Token, error)
}

// Identity represents a simplified identity for drive operations.
type Identity struct {
	ID       string
	Provider string
}

// Token represents a simplified access token.
type Token struct {
	AccessToken string
}
