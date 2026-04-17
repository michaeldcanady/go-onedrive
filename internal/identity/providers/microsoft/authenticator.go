package microsoft

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal/identity/shared"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/state"
)

const (
	tokenBucket = "tokens/microsoft"
)

// Authenticator implements the identity.shared.Authenticator interface for Microsoft.
type Authenticator struct {
	creds map[string]azcore.TokenCredential
	mu    sync.RWMutex
	state state.Service
	log   logger.Logger
}

// NewAuthenticator initializes a new Microsoft authenticator.
func NewAuthenticator(state state.Service, log logger.Logger) *Authenticator {
	return &Authenticator{
		creds: make(map[string]azcore.TokenCredential),
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
	a.log.Info("starting microsoft authentication", logger.String("method", opts.Method.String()), logger.String("identity", opts.IdentityID))

	identityID := opts.IdentityID

	if opts.Force {
		a.mu.Lock()
		delete(a.creds, identityID)
		a.mu.Unlock()
	}

	cred, err := a.getOrUpdateCredential(ctx, identityID, opts)
	if err != nil {
		return shared.AccessToken{}, err
	}

	// Common scopes for OneDrive
	scopes := []string{"https://graph.microsoft.com/.default"}
	token, err := cred.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: scopes,
	})
	if err != nil {
		return shared.AccessToken{}, fmt.Errorf("failed to get token: %w", err)
	}

	// For Microsoft, we might want to get the actual user ID/email from the token if not provided.
	// For now, assume it's provided or we use a placeholder if empty.

	return shared.AccessToken{
		IdentityID: identityID,
		Token:      token.Token,
		ExpiresAt:  token.ExpiresOn,
		Scopes:     scopes,
	}, nil
}

// SaveToken persists the provided access token to state.
func (a *Authenticator) SaveToken(ctx context.Context, token shared.AccessToken) error {
	a.log.Debug("caching access token", logger.String("identity", token.IdentityID))

	tokenData, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to serialize token for caching: %w", err)
	}

	if err := a.state.SetScoped(tokenBucket, token.IdentityID, string(tokenData), state.ScopeGlobal); err != nil {
		return fmt.Errorf("failed to cache access token: %w", err)
	}

	return nil
}

func (a *Authenticator) getOrUpdateCredential(ctx context.Context, identityID string, opts shared.LoginOptions) (azcore.TokenCredential, error) {
	a.mu.RLock()
	cred, ok := a.creds[identityID]
	a.mu.RUnlock()

	if ok && !opts.Force {
		return cred, nil
	}

	newCred, err := a.createCredential(opts)
	if err != nil {
		return nil, err
	}

	a.mu.Lock()
	a.creds[identityID] = newCred
	a.mu.Unlock()

	return newCred, nil
}

func (a *Authenticator) createCredential(opts shared.LoginOptions) (azcore.TokenCredential, error) {
	tenantID := opts.ProviderSpecific["tenant_id"]
	clientID := opts.ProviderSpecific["client_id"]

	switch opts.Method {
	case shared.AuthMethodInteractiveBrowser, shared.AuthMethodUnknown:
		if !opts.Interactive {
			// Try to see if we have a cached token that we can use to create a static credential
			// This is a bit of a shim for the VFS architecture.
			// Ideally we use a Refreshable credential.
			return nil, errors.New("interactive authentication not allowed in non-interactive mode")
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

// Logout removes cached credentials for a specific identity.
func (a *Authenticator) Logout(ctx context.Context, identityID string) error {
	a.log.Info("logging out from microsoft", logger.String("identity", identityID))
	a.mu.Lock()
	delete(a.creds, identityID)
	a.mu.Unlock()
	return a.state.ClearScoped(tokenBucket, identityID)
}

// GetCredential returns the underlying Azure token credential for a specific identity.
func (a *Authenticator) GetCredential(ctx context.Context, identityID string) (any, error) {
	a.mu.RLock()
	cred, ok := a.creds[identityID]
	a.mu.RUnlock()

	if ok {
		return cred, nil
	}

	// Try to load from state
	tokenData, err := a.state.GetScoped(tokenBucket, identityID)
	if err != nil {
		return nil, fmt.Errorf("no credential found for identity %s: %w", identityID, err)
	}

	var token shared.AccessToken
	if err := json.Unmarshal([]byte(tokenData), &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached token: %w", err)
	}

	// For now, return a static token credential.
	// In the future, we should store enough info to recreate a refreshable credential.
	cred = NewStaticTokenCredential(token)

	a.mu.Lock()
	a.creds[identityID] = cred
	a.mu.Unlock()

	return cred, nil
}

// ListIdentities returns all cached Microsoft identity IDs.
func (a *Authenticator) ListIdentities(ctx context.Context) ([]string, error) {
	return a.state.ListScoped(tokenBucket)
}
