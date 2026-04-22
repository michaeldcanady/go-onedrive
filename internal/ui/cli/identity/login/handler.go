package login

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/features/config"
	"github.com/michaeldcanady/go-onedrive/internal/identity"
	"github.com/michaeldcanady/go-onedrive/internal/features/logger"
)

// Command orchestrates the authentication flow for a specific request.
type Command struct {
	config   config.Service
	identity identity.Service
	log      logger.Logger
}

// NewCommand initializes a new instance of the login Command.
func NewCommand(
	cfg config.Service,
	id identity.Service,
	l logger.Logger,
) *Command {
	return &Command{
		config:   cfg,
		identity: id,
		log:      l,
	}
}

// Validate prepares and validates the options for the login operation.
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	return opts.Validate()
}

// Execute executes the login operation based on the current profile and provided options.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx)

	log.Info("starting login flow")

	cfg, err := c.config.GetConfig(ctx)
	if err != nil {
		log.Error("failed to load configuration", logger.Error(err))
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	provider := cfg.Auth.Provider
	if provider == "" {
		provider = "microsoft"
	}

	method := identity.AuthMethodUnknown
	if opts.Method != "" {
		method = identity.ParseAuthMethod(opts.Method)
	} else if cfg.Auth.Method != "" {
		method = identity.ParseAuthMethod(cfg.Auth.Method)
	}

	identityID := opts.IdentityID
	if opts.Alias != "" {
		identityID = opts.Alias
	}

	loginOpts := identity.LoginOptions{
		AccountID:   identityID,
		Force:       opts.Force,
		Interactive: true,
		Method:      method,
		ProviderSpecific: map[string]string{
			"tenant_id":     opts.TenantID,
			"client_id":     opts.ClientID,
			"client_secret": opts.ClientSecret,
		},
	}

	if loginOpts.ProviderSpecific["tenant_id"] == "" {
		loginOpts.ProviderSpecific["tenant_id"] = cfg.Auth.TenantID
	}
	if loginOpts.ProviderSpecific["client_id"] == "" {
		loginOpts.ProviderSpecific["client_id"] = cfg.Auth.ClientID
	}
	if loginOpts.ProviderSpecific["client_secret"] == "" {
		loginOpts.ProviderSpecific["client_secret"] = cfg.Auth.ClientSecret
	}

	log.Info("authenticating", logger.String("provider", provider), logger.String("method", method.String()))
	req, err := identity.ToProtoAuthenticateRequest(loginOpts)
	if err != nil {
		return fmt.Errorf("failed to prepare authentication request: %w", err)
	}

	resp, err := c.identity.Authenticate(ctx, provider, req)
	if err != nil {
		log.Error("authentication failed", logger.Error(err))
		return fmt.Errorf("authentication failed: %w", err)
	}
	token := identity.FromProtoAccessToken(resp.GetToken(), identityID)
	acc := identity.FromProtoIdentity(resp.GetIdentity())

	// Ensure the AccountID is populated from the identity record
	token.AccountID = acc.ID

	if err := c.identity.GetStore().Save(ctx, provider, token); err != nil {
		log.Error("failed to cache token", logger.Error(err))
		return fmt.Errorf("failed to cache access token: %w", err)
	}
	log.Info("authentication successful and token cached",
		logger.String("identity", acc.ID),
		logger.String("display_name", acc.DisplayName),
		logger.String("email", acc.Email))

	if opts.ShowToken {
		fmt.Fprintf(opts.Stdout, "Access Token: %s\n", token.Token)
	}

	return nil
}

// Finalize performs any necessary cleanup after the login operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
