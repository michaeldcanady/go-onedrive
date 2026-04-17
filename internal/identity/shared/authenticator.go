package shared

import (
	"context"
)

// Authenticator defines the interface for an identity provider's authentication logic.
type Authenticator interface {
	// ProviderName returns the unique identifier for this provider (e.g., "microsoft").
	ProviderName() string
	// Authenticate performs the authentication flow and returns the resulting token.
	Authenticate(ctx context.Context, opts LoginOptions) (AccessToken, error)
	// SaveToken persists the provided access token.
	SaveToken(ctx context.Context, token AccessToken) error
	// Logout removes any cached authentication state for a specific identity.
	Logout(ctx context.Context, identityID string) error
	// GetCredential returns a provider-specific credential object (e.g., azcore.TokenCredential).
	GetCredential(ctx context.Context, identityID string) (any, error)
}
