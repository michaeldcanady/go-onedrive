package microsoft

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	interrors "github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/identity"
	"github.com/michaeldcanady/go-onedrive/internal/identity/shared"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Authenticator implements the identity.shared.Authenticator interface for Microsoft.
type Authenticator struct {
	creds map[string]azcore.TokenCredential
	mu    sync.RWMutex
	repo  identity.TokenRepository
	log   logger.Logger
}

// NewAuthenticator initializes a new Microsoft authenticator.
func NewAuthenticator(repo identity.TokenRepository, log logger.Logger) *Authenticator {
	return &Authenticator{
		creds: make(map[string]azcore.TokenCredential),
		repo:  repo,
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
	} else if identityID != "" {
		// If not forcing and identity is known, check if we already have a valid token
		if token, err := a.repo.Get(ctx, a.ProviderName(), identityID); err == nil {
			if token.ExpiresAt.After(time.Now().Add(5 * time.Minute)) {
				a.log.Info("using valid cached token from repository", logger.String("identity", identityID))
				return token, nil
			}
		}
	} else {
		// If identity is unknown and not forcing, check if there's exactly one identity with a valid token
		if identities, err := a.ListIdentities(ctx); err == nil && len(identities) == 1 {
			if token, err := a.repo.Get(ctx, a.ProviderName(), identities[0]); err == nil {
				if token.ExpiresAt.After(time.Now().Add(5 * time.Minute)) {
					a.log.Info("using valid cached token for the only available identity", logger.String("identity", identities[0]))
					return token, nil
				}
			}
		}
	}

	cred, err := a.getOrUpdateCredential(ctx, identityID, opts)
	if err != nil {
		return shared.AccessToken{}, err
	}

	// Common scopes for OneDrive
	scopes := []string{"https://graph.microsoft.com/.default"}

	// If identity is unknown, try to authenticate first to get the identity record.
	// This avoids double prompts because Authenticate will trigger the login flow,
	// and subsequent GetToken calls will use the cached result.
	if identityID == "" {
		type authenticatable interface {
			Authenticate(context.Context, *policy.TokenRequestOptions) (azidentity.AuthenticationRecord, error)
		}

		if aCred, ok := cred.(authenticatable); ok {
			record, err := aCred.Authenticate(ctx, &policy.TokenRequestOptions{
				Scopes: scopes,
			})
			if err == nil && record.Username != "" {
				identityID = record.Username
				a.log.Info("obtained identity from credential record", logger.String("identity", identityID))

				// Update creds map with the new identity ID immediately
				a.mu.Lock()
				a.creds[identityID] = cred
				delete(a.creds, "")
				a.mu.Unlock()
			} else if err != nil {
				a.log.Debug("failed to get authentication record", logger.Error(err))
			}
		}
	}

	token, err := cred.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: scopes,
	})
	if err != nil {
		return shared.AccessToken{}, fmt.Errorf("failed to get token: %w", err)
	}

	if identityID == "" {
		// Fallback to extracting from token if still empty (e.g. for non-authenticatable creds)
		extractedID, err := extractIdentityFromToken(token.Token)
		if err != nil {
			a.log.Debug("failed to extract identity from token", logger.Error(err))
			identityID = "unknown"
		} else {
			identityID = extractedID
			a.log.Info("extracted identity from token", logger.String("identity", identityID))
		}

		// Update creds map with the discovered identity ID
		if identityID != "" && identityID != "unknown" {
			a.mu.Lock()
			a.creds[identityID] = cred
			delete(a.creds, "")
			a.mu.Unlock()
		}
	}

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

	if err := a.repo.Save(ctx, a.ProviderName(), token); err != nil {
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
			return nil, fmt.Errorf("interactive authentication not allowed in non-interactive mode")
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
	return a.repo.Delete(ctx, a.ProviderName(), identityID)
}

// GetCredential returns the underlying Azure token credential for a specific identity.
func (a *Authenticator) GetCredential(ctx context.Context, identityID string) (any, error) {
	if identityID == "" {
		identities, err := a.ListIdentities(ctx)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to list identities: %w", interrors.ErrUnauthorized, err)
		}

		if len(identities) == 0 {
			return nil, fmt.Errorf("%w: no identities found, please login using 'odc auth login'", interrors.ErrUnauthorized)
		}

		if len(identities) > 1 {
			return nil, fmt.Errorf("%w: multiple identities found, please specify one with --identity", interrors.ErrUnauthorized)
		}

		identityID = identities[0]
		a.log.Debug("defaulting to the only available identity", logger.String("identity", identityID))
	}

	a.mu.RLock()
	cred, ok := a.creds[identityID]
	a.mu.RUnlock()

	if ok {
		return cred, nil
	}

	// Try to load from repo
	token, err := a.repo.Get(ctx, a.ProviderName(), identityID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to load cached token for identity %s: %w", interrors.ErrUnauthorized, identityID, err)
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
	return a.repo.List(ctx, a.ProviderName())
}
