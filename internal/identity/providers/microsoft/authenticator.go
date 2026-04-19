package microsoft

import (
	"context"
	"fmt"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal/identity"
	proto "github.com/michaeldcanady/go-onedrive/internal/identity/proto"
)

// MicrosoftAuthenticator implements the identity.Authenticator interface for Microsoft.
type MicrosoftAuthenticator struct {
	creds map[string]azcore.TokenCredential
	mu    sync.RWMutex
}

// NewMicrosoftAuthenticator initializes a new Microsoft authenticator.
func NewMicrosoftAuthenticator() *MicrosoftAuthenticator {
	return &MicrosoftAuthenticator{
		creds: make(map[string]azcore.TokenCredential),
	}
}

// ProviderName returns "microsoft".
func (a *MicrosoftAuthenticator) ProviderName() string {
	return "microsoft"
}

// Authenticate performs the Microsoft-specific login flow.
func (a *MicrosoftAuthenticator) Authenticate(ctx context.Context, req *proto.AuthenticateRequest) (*proto.AuthenticateResponse, error) {
	opts, err := identity.FromProtoAuthenticateRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to convert proto request: %w", err)
	}

	accountID := opts.AccountID
	var acc identity.Account

	if opts.Force {
		a.mu.Lock()
		delete(a.creds, accountID)
		a.mu.Unlock()
	}

	// If account ID is unknown, try to authenticate first to get the identity record.
	if accountID == "" {
		type authenticatable interface {
			Authenticate(context.Context, *policy.TokenRequestOptions) (azidentity.AuthenticationRecord, error)
		}

		cred, err := a.getOrUpdateCredential(ctx, accountID, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to get credential: %w", err)
		}

		if aCred, ok := cred.(authenticatable); ok {
			scopes := []string{"https://graph.microsoft.com/.default"} // Example scope
			record, err := aCred.Authenticate(ctx, &policy.TokenRequestOptions{
				Scopes: scopes,
			})
			if err == nil && record.Username != "" {
				accountID = record.Username
				acc.ID = record.HomeAccountID
				acc.Email = record.Username
				acc.DisplayName = record.Username // Fallback to username for display name
				acc.Provider = a.ProviderName()

				// Update creds map with the new identity ID immediately
				a.mu.Lock()
				a.creds[accountID] = cred
				delete(a.creds, "") // Remove placeholder if it existed
				a.mu.Unlock()
			}
		}
	}

	// If accountID is still empty after the above, it means authentication failed or no identity was found.
	if accountID == "" {
		return nil, fmt.Errorf("authentication failed: could not determine account ID")
	}

	// Populate a minimal account if not fully determined
	if acc.ID == "" {
		acc.ID = accountID
		acc.Email = accountID
		acc.DisplayName = accountID
		acc.Provider = a.ProviderName()
	}

	return &proto.AuthenticateResponse{
		Token: &proto.AccessToken{
			Token:        "",
			RefreshToken: "",
			ExpiresAt:    0,
			Scopes:       nil,
		},
		Identity: identity.ToProtoIdentity(acc),
	}, nil
}

// getOrUpdateCredential and createCredential are helper methods for managing credentials internally.
// They are kept here for now, as Authenticate still needs to create credentials during the initial flow.
// This logic might be moved closer to Authorizer logic in a future refactor.
func (a *MicrosoftAuthenticator) getOrUpdateCredential(ctx context.Context, accountID string, opts identity.LoginOptions) (azcore.TokenCredential, error) {
	a.mu.RLock()
	cred, ok := a.creds[accountID]
	a.mu.RUnlock()

	if ok && !opts.Force {
		return cred, nil
	}

	newCred, err := a.createCredential(opts)
	if err != nil {
		return nil, err
	}

	a.mu.Lock()
	a.creds[accountID] = newCred
	a.mu.Unlock()

	return newCred, nil
}

func (a *MicrosoftAuthenticator) createCredential(opts identity.LoginOptions) (azcore.TokenCredential, error) {
	tenantID := opts.ProviderSpecific["tenant_id"]
	clientID := opts.ProviderSpecific["client_id"]

	switch opts.Method {
	case identity.AuthMethodInteractiveBrowser, identity.AuthMethodUnknown:
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

	case identity.AuthMethodDeviceCode:
		return azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
			TenantID: tenantID,
			ClientID: clientID,
		})

	case identity.AuthMethodClientSecret:
		clientSecret := opts.ProviderSpecific["client_secret"]
		if tenantID == "" || clientID == "" || clientSecret == "" {
			return nil, fmt.Errorf("tenant_id, client_id, and client_secret are required for client-secret auth")
		}
		return azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)

	default:
		return nil, fmt.Errorf("unsupported authentication method: %v", opts.Method)
	}
}
