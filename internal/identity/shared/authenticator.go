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
	// Logout removes any cached authentication state for this provider.
	Logout(ctx context.Context) error
}
