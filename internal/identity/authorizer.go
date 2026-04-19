package identity

import (
	"context"

	proto "github.com/michaeldcanady/go-onedrive/internal/identity/proto"
)

// Authorizer provides a way to obtain access tokens for specific accounts, handling refreshes automatically.
type Authorizer interface {
	// Token returns a valid access token for the given identity request.
	Token(ctx context.Context, req *proto.GetTokenRequest) (*proto.GetTokenResponse, error)
}
