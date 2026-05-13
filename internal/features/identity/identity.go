package identity

import (
	"context"
	"time"

	identity_proto "github.com/michaeldcanady/go-onedrive/internal/features/plugins/proto/identity"
)

// Identity represents an authenticated user profile from a specific provider.
type Identity struct {
	ID          string            `json:"id"`
	DisplayName string            `json:"display_name"`
	Email       string            `json:"email"`
	Provider    string            `json:"provider"`
	Metadata    map[string]string `json:"metadata"`
}

// Token encapsulates the credentials required to access protected resources.
type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	Scopes       []string  `json:"scopes"`
}

// Service coordinates the authentication lifecycle, including interactive login
// flows and identity management.
type Service interface {
	// Login initiates an interactive authentication flow through a provider-specific plugin.
	Login(ctx context.Context, provider string, options map[string]string) (*Identity, error)

	// Logout invalidates the session for the specified identity.
	Logout(ctx context.Context, identityID string) error

	// List returns all authenticated identities stored in the local cache.
	List(ctx context.Context) ([]*Identity, error)

	// GetIdentity retrieves a single identity by its unique identifier.
	GetIdentity(ctx context.Context, identityID string) (*Identity, error)

	// FindIdentity searches for an identity matching the query (ID, Email, or DisplayName).
	FindIdentity(ctx context.Context, query string) (*Identity, error)
}

// TokenService manages the lifecycle of authentication tokens, including
// secure storage and proactive background refreshing.
type TokenService interface {
	// GetToken retrieves a valid access token, refreshing it if it is near expiration.
	GetToken(ctx context.Context, provider string, identityID string) (*Token, error)

	// SaveToken persists the provided token for an identity.
	SaveToken(ctx context.Context, provider string, identityID string, token *Token) error

	// RefreshToken forces a token refresh using the stored refresh token.
	RefreshToken(ctx context.Context, provider string, identityID string) (*Token, error)
}

// Repository handles the persistent storage of identities and their associated tokens.
type Repository interface {
	SaveIdentity(i *Identity) error
	GetIdentity(id string) (*Identity, error)
	ListIdentities() ([]*Identity, error)
	DeleteIdentity(id string) error

	SaveToken(provider string, identityID string, t *Token) error
	GetToken(provider string, identityID string) (*Token, error)
	DeleteToken(provider string, identityID string) error
}

func FromProtoIdentity(p *identity_proto.Identity) *Identity {
	return &Identity{
		ID:          p.Id,
		DisplayName: p.DisplayName,
		Email:       p.Email,
		Provider:    p.Provider,
		Metadata:    p.Metadata,
	}
}

func FromProtoToken(p *identity_proto.AccessToken) *Token {
	return &Token{
		AccessToken:  p.AccessToken,
		RefreshToken: p.RefreshToken,
		ExpiresAt:    time.Unix(p.ExpiresAt, 0),
		Scopes:       p.Scopes,
	}
}
