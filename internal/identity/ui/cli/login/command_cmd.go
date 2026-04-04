package login

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/identity/registry"
	"github.com/michaeldcanady/go-onedrive/internal/identity/shared"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Handler orchestrates the authentication flow for a specific request.
type Handler struct {
	config   config.Service
	identity registry.Service
	log      logger.Logger
}

// NewHandler initializes a new instance of the login Handler.
func NewHandler(
	cfg config.Service,
	id registry.Service,
	l logger.Logger,
) *Handler {
	return &Handler{
		config:   cfg,
		identity: id,
		log:      l,
	}
}

// Handle executes the login operation based on the current profile and provided options.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx)

	log.Info("starting login flow")

	log.Debug("loading profile configuration")
	cfg, err := h.config.GetConfig(ctx)
	if err != nil {
		log.Error("failed to load configuration", logger.Error(err))
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	provider := cfg.Auth.Provider
	if provider == "" {
		log.Debug("no provider specified, defaulting to microsoft")
		provider = "microsoft"
	}

	log.Debug("retrieving authenticator for provider", logger.String("provider", provider))
	auth, err := h.identity.Get(provider)
	if err != nil {
		log.Error("unsupported provider", logger.String("provider", provider), logger.Error(err))
		return fmt.Errorf("provider %s not supported: %w", provider, err)
	}

	// Determine method: CLI flag takes precedence, then config
	method := shared.AuthMethodUnknown
	if opts.Method != "" {
		method = shared.ParseAuthMethod(opts.Method)
		log.Debug("auth method set via CLI flag", logger.String("method", method.String()))
	} else if cfg.Auth.Method != "" {
		method = shared.ParseAuthMethod(cfg.Auth.Method)
		log.Debug("auth method set via configuration", logger.String("method", method.String()))
	}

	loginOpts := shared.LoginOptions{
		Force:       true, // Login command always forces a fresh flow
		Interactive: true,
		Method:      method,
		ProviderSpecific: map[string]string{
			"tenant_id":     opts.TenantID,
			"client_id":     opts.ClientID,
			"client_secret": opts.ClientSecret,
		},
	}

	// Merge from config if not provided in CLI
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
	token, err := auth.Authenticate(ctx, loginOpts)
	if err != nil {
		log.Error("authentication failed", logger.Error(err))
		return fmt.Errorf("authentication failed: %w", err)
	}

	if err := auth.SaveToken(ctx, token); err != nil {
		log.Error("failed to cache token", logger.Error(err))
		return fmt.Errorf("failed to cache access token: %w", err)
	}

	log.Info("authentication successful and token cached")

	if opts.ShowToken {
		fmt.Fprintf(opts.Stdout, "Access Token: %s\n", token.Token)
	}

	return nil
}
