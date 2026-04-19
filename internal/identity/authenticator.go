package identity

import (
	"context"

	proto "github.com/michaeldcanady/go-onedrive/internal/identity/proto"
)

// Authenticator defines the interface for an identity provider's authentication logic.
type Authenticator interface {
	// ProviderName returns the unique identifier for this provider (e.g., "microsoft").
	ProviderName() string
	// Authenticate performs the authentication flow and returns the resulting token and identity metadata.
	Authenticate(ctx context.Context, req *proto.AuthenticateRequest) (*proto.AuthenticateResponse, error)
}
