package microsoft

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/identity"
	proto "github.com/michaeldcanady/go-onedrive/internal/identity/proto"
)

// MicrosoftAuthorizer implements the identity.Authorizer interface for Microsoft.
type MicrosoftAuthorizer struct {
	store identity.AccountStore
}

// NewMicrosoftAuthorizer initializes a new Microsoft authorizer.
func NewMicrosoftAuthorizer(store identity.AccountStore) *MicrosoftAuthorizer {
	return &MicrosoftAuthorizer{
		store: store,
	}
}

// Token returns a valid access token for the given identity request.
func (a *MicrosoftAuthorizer) Token(ctx context.Context, req *proto.GetTokenRequest) (*proto.GetTokenResponse, error) {
	accountID := req.GetIdentityId()

	// Get the stored access token (which is already a credential-like object in our architecture)
	token, err := a.store.Get(ctx, "microsoft", accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get token from store: %w", err)
	}

	accessToken := identity.AccessToken{
		AccountID: accountID,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		Scopes:    token.Scopes,
	}
	return &proto.GetTokenResponse{
		Token: identity.ToProtoAccessToken(accessToken),
	}, nil
}
