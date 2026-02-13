package msalbuilder

import (
	"context"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	domainauth "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	msalclient "github.com/michaeldcanady/go-onedrive/internal2/infra/auth/msal/client"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
)

type ClientSecretBuilder struct {
	factory  msalclient.Factory
	cfg      config.AuthenticationConfigImpl
	opts     []confidential.AcquireByCredentialOption
	flowOpts flowOptions
}

func (f *ClientSecretBuilder) addOpt(opt confidential.AcquireByCredentialOption) {
	f.opts = append(f.opts, opt)
}

func (f *ClientSecretBuilder) WithClaims(claims string) *ClientSecretBuilder {
	f.addOpt(confidential.WithClaims(claims))
	return f
}

func (f *ClientSecretBuilder) WithTenantID(tenantID string) *ClientSecretBuilder {
	f.addOpt(confidential.WithTenantID(tenantID))
	return f
}

func (f *ClientSecretBuilder) WithFMIPath(path string) *ClientSecretBuilder {
	f.addOpt(confidential.WithFMIPath(path))
	return f
}
func (f *ClientSecretBuilder) WithAttribute(attrValue string) *ClientSecretBuilder {
	f.addOpt(confidential.WithAttribute(attrValue))
	return f
}

func (f *ClientSecretBuilder) Acquire(ctx context.Context) (domainauth.TokenResult, error) {
	client, err := f.factory.ConfidentialClient(f.cfg)
	if err != nil {
		return domainauth.TokenResult{}, err
	}

	res, err := client.AcquireTokenByCredential(ctx, f.flowOpts.scopes)
	if err != nil {
		return domainauth.TokenResult{}, err
	}

	return domainauth.TokenResult{
		AccessToken: res.AccessToken,
		ExpiresOn:   res.ExpiresOn,
	}, nil
}
