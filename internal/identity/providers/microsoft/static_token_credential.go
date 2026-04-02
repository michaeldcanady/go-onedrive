package microsoft

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/michaeldcanady/go-onedrive/internal/identity/shared"
)

// StaticTokenCredential implements azcore.TokenCredential using a static access token.
// It returns an error if the token is expired.
type StaticTokenCredential struct {
	token shared.AccessToken
}

// NewStaticTokenCredential creates a new StaticTokenCredential.
func NewStaticTokenCredential(token shared.AccessToken) *StaticTokenCredential {
	return &StaticTokenCredential{
		token: token,
	}
}

// GetToken checks for token expiry and returns the token if it's valid.
func (c *StaticTokenCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	if c.isExpired() {
		return azcore.AccessToken{}, fmt.Errorf("cached access token is expired, please run 'login' again")
	}

	return azcore.AccessToken{
		Token:     c.token.Token,
		ExpiresOn: c.token.ExpiresAt,
	}, nil
}

func (c *StaticTokenCredential) isExpired() bool {
	// Check if the token is expired, with a small buffer (e.g., 5 minutes)
	return c.token.ExpiresAt.Before(time.Now().Add(5 * time.Minute))
}
