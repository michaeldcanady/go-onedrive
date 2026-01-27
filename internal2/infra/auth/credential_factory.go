package auth

import (
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
)

type CredentialFactory interface {
	Create(authCfg config.AuthenticationConfig) (azcore.TokenCredential, error)
}

type DefaultCredentialFactory struct{}

func NewDefaultCredentialFactory() *DefaultCredentialFactory {
	return &DefaultCredentialFactory{}
}

func (f *DefaultCredentialFactory) Create(authCfg config.AuthenticationConfig) (azcore.TokenCredential, error) {
	switch authCfg.GetAuthenticationType() {
	case InteractiveBrowserAuthType:
		return f.createInteractiveBrowser(authCfg)
	default:
		return nil, fmt.Errorf("unsupported authentication type: %s", authCfg.GetAuthenticationType())
	}
}

func (f *DefaultCredentialFactory) createInteractiveBrowser(authCfg config.AuthenticationConfig) (azcore.TokenCredential, error) {
	var record azidentity.AuthenticationRecord
	if r := authCfg.GetAuthenticationRecord(); r != nil {
		record = *r
	}

	tokenCache, err := cache.New(nil)
	if err != nil {
		return nil, errors.Join(errors.New("failed to initialize token cache"), err)
	}

	return azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
		TenantID:                       authCfg.GetTenantID(),
		ClientID:                       authCfg.GetClientID(),
		RedirectURL:                    authCfg.GetRedirectURI(),
		DisableAutomaticAuthentication: true,
		AuthenticationRecord:           record,
		Cache:                          tokenCache,
	})
}
