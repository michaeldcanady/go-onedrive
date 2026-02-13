package msalbuilder

import (
	"context"

	domainauth "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	msalclient "github.com/michaeldcanady/go-onedrive/internal2/infra/auth/msal/client"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
)

type Flow interface {
	Acquire(ctx context.Context) (domainauth.TokenResult, error)
}

type flowOptions struct {
	scopes []string
}

type FlowBuilder struct {
	factory msalclient.Factory
	cfg     config.AuthenticationConfigImpl
	opts    flowOptions
}

func NewFlowBuilder(factory msalclient.Factory, cfg config.AuthenticationConfigImpl) *FlowBuilder {
	return &FlowBuilder{factory: factory, cfg: cfg}
}

func (b *FlowBuilder) WithScopes(scopes []string) *FlowBuilder {
	b.opts.scopes = scopes
	return b
}

func (b *FlowBuilder) DeviceCode() *DeviceCodeBuilder {
	return &DeviceCodeBuilder{factory: b.factory, cfg: b.cfg, opts: b.opts}
}

func (b *FlowBuilder) Interactive() *InteractiveBuilder {
	return &InteractiveBuilder{factory: b.factory, cfg: b.cfg, opts: b.opts}
}

func (b *FlowBuilder) ROPC() *ROPCBuilder {
	return &ROPCBuilder{factory: b.factory, cfg: b.cfg, opts: b.opts}
}

func (b *FlowBuilder) ClientSecret() *ClientSecretBuilder {
	return &ClientSecretBuilder{factory: b.factory, cfg: b.cfg, flowOpts: b.opts}
}
