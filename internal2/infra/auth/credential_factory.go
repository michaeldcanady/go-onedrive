package auth

import (
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/account"
)

type CredentialFactory interface {
	Create(authCfg CredentialOptions) (azcore.TokenCredential, error)
}

type DefaultCredentialFactory struct{}

func NewDefaultCredentialFactory() *DefaultCredentialFactory {
	return &DefaultCredentialFactory{}
}

func (f *DefaultCredentialFactory) Create(authCfg CredentialOptions) (azcore.TokenCredential, error) {
	switch authCfg.Type {
	case InteractiveBrowserAuthType:
		return f.createInteractiveBrowser(authCfg)
	default:
		return nil, fmt.Errorf("unsupported authentication type: %s", authCfg.Type)
	}
}

func (f *DefaultCredentialFactory) createInteractiveBrowser(authCfg CredentialOptions) (azcore.TokenCredential, error) {
	tokenCache, err := cache.New(nil)
	if err != nil {
		return nil, errors.Join(errors.New("failed to initialize token cache"), err)
	}

	var authenticationRecord azidentity.AuthenticationRecord
	if authCfg.Account != (account.Account{}) {
		authenticationRecord = account.AccountToMSAuthRecord(authCfg.Account)
	} else if authCfg.AuthenticationRecord != (azidentity.AuthenticationRecord{}) {
		authenticationRecord = authCfg.AuthenticationRecord
	}

	return azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
		TenantID:                       authCfg.TenantID,
		ClientID:                       authCfg.ClientID,
		RedirectURL:                    authCfg.RedirectURI,
		DisableAutomaticAuthentication: true,
		AuthenticationRecord:           authenticationRecord,
		Cache:                          tokenCache,
	})
}
