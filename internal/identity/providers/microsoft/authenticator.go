package microsoft

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal/identity/shared"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/state"
)

// Authenticator implements the identity.shared.Authenticator interface for Microsoft.
type Authenticator struct {
	cred  azcore.TokenCredential
	state state.AuthStore
	log   logger.Logger
}

// NewAuthenticator initializes a new Microsoft authenticator.
func NewAuthenticator(cred azcore.TokenCredential, state state.AuthStore, log logger.Logger) *Authenticator {
	return &Authenticator{
		cred:  cred,
		state: state,
		log:   log,
	}
}

// ProviderName returns "microsoft".
func (a *Authenticator) ProviderName() string {
	return "microsoft"
}

// Authenticate performs the Microsoft-specific login flow.
func (a *Authenticator) Authenticate(ctx context.Context, opts shared.LoginOptions) (shared.AccessToken, error) {
	a.log.Info("starting microsoft authentication", logger.String("method", opts.Method.String()))

	if opts.Force {
		a.cred = nil
	}

	if a.cred == nil {
		cred, err := a.createCredential(opts)
		if err != nil {
			return shared.AccessToken{}, err
		}
		a.cred = cred
	}

	// Common scopes for OneDrive
	scopes := []string{"https://graph.microsoft.com/.default"}
	token, err := a.cred.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: scopes,
	})
	if err != nil {
		return shared.AccessToken{}, fmt.Errorf("failed to get token: %w", err)
	}

	return shared.AccessToken{
		Token:     token.Token,
		ExpiresAt: token.ExpiresOn,
		Scopes:    scopes,
	}, nil
}

// SaveToken persists the provided access token to state.
func (a *Authenticator) SaveToken(ctx context.Context, token shared.AccessToken) error {
	a.log.Debug("caching access token")

	tokenData, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to serialize token for caching: %w", err)
	}

	if err := a.state.SetAccessToken(ctx, string(tokenData), state.ScopeGlobal); err != nil {
		return fmt.Errorf("failed to cache access token: %w", err)
	}

	return nil
}

func (a *Authenticator) createCredential(opts shared.LoginOptions) (azcore.TokenCredential, error) {
	tenantID := opts.ProviderSpecific["tenant_id"]
	clientID := opts.ProviderSpecific["client_id"]

	switch opts.Method {
	case shared.AuthMethodInteractiveBrowser, shared.AuthMethodUnknown:
		if !opts.Interactive {
			return nil, errors.New("interactive authentication not allowed")
		}

		return azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
			TenantID: tenantID,
			ClientID: clientID,
		})

	case shared.AuthMethodDeviceCode:

		return azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
			TenantID: tenantID,
			ClientID: clientID,
		})

	case shared.AuthMethodClientSecret:
		clientSecret := opts.ProviderSpecific["client_secret"]
		if tenantID == "" || clientID == "" || clientSecret == "" {
			return nil, fmt.Errorf("tenant_id, client_id, and client_secret are required for client-secret auth")
		}
		return azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)

	default:
		return nil, fmt.Errorf("unsupported authentication method: %v", opts.Method)
	}
}

// Logout removes cached credentials.
func (a *Authenticator) Logout(ctx context.Context) error {
	a.log.Info("logging out from microsoft")
	a.cred = nil
	return a.state.ClearAccessToken(ctx)
}

// Credential returns the underlying Azure token credential.
func (a *Authenticator) Credential() azcore.TokenCredential {
	return a.cred
}

