package microsoft

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/michaeldcanady/go-onedrive/internal/identity"
	proto "github.com/michaeldcanady/go-onedrive/internal/identity/proto"
)

// MicrosoftAuthorizer implements the identity.Authorizer interface for Microsoft.
type MicrosoftAuthorizer struct {
	cred azcore.TokenCredential
}

// NewMicrosoftAuthorizer initializes a new Microsoft authorizer.
func NewMicrosoftAuthorizer(cred azcore.TokenCredential) *MicrosoftAuthorizer {
	return &MicrosoftAuthorizer{
		cred: cred,
	}
}

// Token returns a valid access token for the given identity request.
func (a *MicrosoftAuthorizer) Token(ctx context.Context, req *proto.GetTokenRequest) (*proto.GetTokenResponse, error) {
	accountID := req.GetIdentityId()

	// Get a new token from the credential
	scopes := req.GetScopes()
	if len(scopes) == 0 {
		scopes = []string{"https://graph.microsoft.com/.default"}
	}

	credToken, err := a.cred.GetToken(ctx, policy.TokenRequestOptions{Scopes: scopes})
	if err != nil {
		return nil, fmt.Errorf("failed to get token from credential: %w", err)
	}
	// Try to extract richer identity info from the token
	_, err = extractFullIdentityFromToken(credToken.Token)
	if err != nil {
		// Fallback to a minimal identity if extraction fails
		// Log/handle fallback as needed
	}
	accessToken := identity.AccessToken{
		AccountID: accountID,
		Token:     credToken.Token,
		ExpiresAt: credToken.ExpiresOn,
		Scopes:    scopes,
	}
	return &proto.GetTokenResponse{
		Token: identity.ToProtoAccessToken(accessToken),
	}, nil
}
