package login

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal2/core/config"
	"github.com/michaeldcanady/go-onedrive/internal2/core/identity/registry"
	"github.com/michaeldcanady/go-onedrive/internal2/core/identity/shared"
	"github.com/michaeldcanady/go-onedrive/internal2/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/core/state"
)

// Handler orchestrates the authentication flow for a specific request.
type Handler struct {
	config   config.Service
	state    state.Service
	identity registry.Service
	log      logger.Logger
}

// NewHandler initializes a new instance of the login Handler.
func NewHandler(
	cfg config.Service,
	st state.Service,
	id registry.Service,
	l logger.Logger,
) *Handler {
	return &Handler{
		config:   cfg,
		state:    st,
		identity: id,
		log:      l,
	}
}

// Handle executes the login operation based on the current profile and provided options.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("starting login flow")

	profile, err := h.state.Get(state.KeyProfile)
	if err != nil {
		return fmt.Errorf("failed to get current profile: %w", err)
	}

	cfg, err := h.config.GetConfig(ctx, profile)
	if err != nil {
		return fmt.Errorf("failed to load configuration for profile %s: %w", profile, err)
	}

	provider := cfg.Auth.Provider
	if provider == "" {
		provider = "microsoft" // Default fallback
	}

	auth, err := h.identity.Get(provider)
	if err != nil {
		return fmt.Errorf("provider %s not supported: %w", provider, err)
	}

	h.log.Info("authenticating with provider", logger.String("provider", provider))

	// Determine method: CLI flag takes precedence, then config
	method := shared.AuthMethodUnknown
	if opts.Method != "" {
		method = shared.ParseAuthMethod(opts.Method)
	} else if cfg.Auth.Method != "" {
		// We might need to update the core/config to support Method field if it's not there yet.
		// For now assume it's a string in the config.
		method = shared.ParseAuthMethod(cfg.Auth.Method)
	}

	loginOpts := shared.LoginOptions{
		Force:       opts.Force,
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

	token, err := auth.Authenticate(ctx, loginOpts)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	h.log.Info("authentication successful", logger.String("profile", profile))

	if opts.ShowToken {
		fmt.Fprintf(opts.Stdout, "Access Token: %s\n", token.Token)
	}

	return nil
}
