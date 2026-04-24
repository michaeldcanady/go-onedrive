package identity

import (
	"context"
)

// AccountStore defines the persistence interface for access tokens.
type AccountStore interface {
	// Get retrieves a token for a specific provider and identity.
	Get(ctx context.Context, provider, identityID string) (AccessToken, error)
	// Save persists a token.
	Save(ctx context.Context, provider string, token AccessToken) error
	// Delete removes a token.
	Delete(ctx context.Context, provider, identityID string) error
	// List returns all identity IDs for a specific provider.
	List(ctx context.Context, provider string) ([]string, error)
}
