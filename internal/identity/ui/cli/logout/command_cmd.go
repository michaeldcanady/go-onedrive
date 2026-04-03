package logout

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/identity/registry"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/state"
)

// Handler orchestrates the logout flow for the current profile.
type Handler struct {
	config   config.Service
	state    state.Service
	identity registry.Service
	log      logger.Logger
}

// NewHandler initializes a new logout Handler.
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

// Handle clears the authentication state for the current profile.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("starting logout flow")

	cfg, err := h.config.GetConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	provider := cfg.Auth.Provider
	if provider == "" {
		provider = "microsoft"
	}

	auth, err := h.identity.Get(provider)
	if err != nil {
		return fmt.Errorf("provider %s not supported: %w", provider, err)
	}

	h.log.Info("performing provider logout", logger.String("provider", provider))
	if err := auth.Logout(ctx); err != nil {
		h.log.Warn("provider logout failed", logger.Error(err))
	}

	h.log.Info("logout successful")
	fmt.Fprintf(opts.Stdout, "Logged out successfully\n")

	return nil
}
