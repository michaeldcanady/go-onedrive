package msalbuilder

import (
	"context"

	domainauth "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	msalclient "github.com/michaeldcanady/go-onedrive/internal2/infra/auth/msal/client"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
)

// DeviceCodeMessage contains the information a user needs to complete authentication.
type DeviceCodeMessage struct {
	// UserCode is the user code returned by the service.
	UserCode string `json:"user_code"`
	// VerificationURL is the URL at which the user must authenticate.
	VerificationURL string `json:"verification_uri"`
	// Message is user instruction from Microsoft Entra ID.
	Message string `json:"message"`
}

type DeviceCodeBuilder struct {
	factory msalclient.Factory
	cfg     config.AuthenticationConfigImpl
	opts    flowOptions
	onMsg   func(context.Context, DeviceCodeMessage) error
}

func (b *DeviceCodeBuilder) WithMessageHandler(fn func(context.Context, DeviceCodeMessage) error) *DeviceCodeBuilder {
	b.onMsg = fn
	return b
}

func (f *DeviceCodeBuilder) Acquire(ctx context.Context) (domainauth.TokenResult, error) {
	client, err := f.factory.PublicClient(f.cfg)
	if err != nil {
		return domainauth.TokenResult{}, err
	}

	dc, err := client.AcquireTokenByDeviceCode(ctx, f.opts.scopes)
	if err != nil {
		return domainauth.TokenResult{}, err
	}

	fn := f.onMsg
	if fn == nil {
		fn = func(ctx context.Context, dm DeviceCodeMessage) error {
			if err := ctx.Err(); err != nil {
				return err
			}
			println(dm.Message)
			return nil
		}
	}

	if err := f.onMsg(ctx, DeviceCodeMessage{
		UserCode:        dc.Result.UserCode,
		VerificationURL: dc.Result.VerificationURL,
		Message:         dc.Result.Message,
	}); err != nil {
		return domainauth.TokenResult{}, err
	}

	res, err := dc.AuthenticationResult(ctx)
	if err != nil {
		return domainauth.TokenResult{}, err
	}

	return domainauth.TokenResult{
		AccessToken: res.AccessToken,
		ExpiresOn:   res.ExpiresOn,
		Account:     res.Account,
	}, nil
}
