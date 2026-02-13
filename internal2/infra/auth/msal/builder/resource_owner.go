package msalbuilder

import (
	"context"

	domainauth "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	msalclient "github.com/michaeldcanady/go-onedrive/internal2/infra/auth/msal/client"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
)

type ROPCBuilder struct {
	factory  msalclient.Factory
	cfg      config.AuthenticationConfigImpl
	opts     flowOptions
	username string
	password string
}

func (b *ROPCBuilder) WithCredentials(u, p string) *ROPCBuilder {
	b.username = u
	b.password = p
	return b
}

func (f *ROPCBuilder) Acquire(ctx context.Context) (domainauth.TokenResult, error) {
	client, err := f.factory.PublicClient(f.cfg)
	if err != nil {
		return domainauth.TokenResult{}, err
	}

	res, err := client.AcquireTokenByUsernamePassword(ctx, f.opts.scopes, f.username, f.password)
	if err != nil {
		return domainauth.TokenResult{}, err
	}

	return domainauth.TokenResult{
		AccessToken: res.AccessToken,
		ExpiresOn:   res.ExpiresOn,
		Account:     res.Account,
	}, nil
}
