package msalbuilder

import (
	"context"

	domainauth "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	msalclient "github.com/michaeldcanady/go-onedrive/internal2/infra/auth/msal/client"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
)

type InteractiveBuilder struct {
	factory msalclient.Factory
	cfg     config.AuthenticationConfigImpl
	opts    flowOptions
}

func (f *InteractiveBuilder) Acquire(ctx context.Context) (domainauth.TokenResult, error) {
	client, err := f.factory.PublicClient(f.cfg)
	if err != nil {
		return domainauth.TokenResult{}, err
	}

	res, err := client.AcquireTokenInteractive(ctx, f.opts.scopes)
	if err != nil {
		return domainauth.TokenResult{}, err
	}

	return domainauth.TokenResult{
		AccessToken: res.AccessToken,
		ExpiresOn:   res.ExpiresOn,
		Account:     res.Account,
	}, nil
}
