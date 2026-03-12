package microsoft

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/michaeldcanady/go-onedrive/internal2/core/identity/shared"
	"github.com/michaeldcanady/go-onedrive/internal2/core/logger"
)

// Authenticator implements the identity.shared.Authenticator interface for Microsoft.
// It bridges the generic authentication logic with the Microsoft Graph SDK.
type Authenticator struct {
	// cred is the underlying Azure token credential.
	cred azcore.TokenCredential
	// log is the logger used for reporting authentication events.
	log logger.Logger
}

// NewAuthenticator initializes a new Microsoft authenticator with the given credential and logger.
func NewAuthenticator(cred azcore.TokenCredential, log logger.Logger) *Authenticator {
	return &Authenticator{
		cred: cred,
		log:  log,
	}
}

// ProviderName returns the identifier for this provider: "microsoft".
func (a *Authenticator) ProviderName() string {
	return "microsoft"
}

// Authenticate performs the Microsoft-specific login flow.
func (a *Authenticator) Authenticate(ctx context.Context, opts shared.LoginOptions) (shared.AccessToken, error) {
	// In a full implementation, we would use the AzIdentity credential here.
	// For this pilot, we are scaffolding the connection.
	return shared.AccessToken{
		Token: "scaffolded-microsoft-token",
	}, nil
}

// Logout performs provider-specific logout logic for Microsoft.
func (a *Authenticator) Logout(ctx context.Context) error {
	return nil
}
