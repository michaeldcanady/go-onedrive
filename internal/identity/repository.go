package identity

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/identity/shared"
)

// TokenRepository defines the persistence interface for access tokens.
type TokenRepository interface {
	// Get retrieves a token for a specific provider and identity.
	Get(ctx context.Context, provider, identityID string) (shared.AccessToken, error)
	// Save persists a token.
	Save(ctx context.Context, provider string, token shared.AccessToken) error
	// Delete removes a token.
	Delete(ctx context.Context, provider, identityID string) error
	// List returns all identity IDs for a specific provider.
	List(ctx context.Context, provider string) ([]string, error)
}
