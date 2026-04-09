package logout

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/identity/registry"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Handler orchestrates the logout flow for the current profile.
type Handler struct {
	config   config.Service
	identity registry.Service
	log      logger.Logger
}

// NewHandler initializes a new logout Handler.
func NewHandler(
	cfg config.Service,
	id registry.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("identity-logout")
	return &Handler{
		config:   cfg,
		identity: id,
		log:      cliLog,
	}
}

// Handle clears the authentication state for the current profile.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx)

	log.Info("starting logout flow")

	cfg, err := h.config.GetConfig(ctx)
	if err != nil {
		log.Error(err.Error(), errors.LogFields(err)...)
		return err
	}

	provider := cfg.Auth.Provider
	if provider == config.AuthProviderUnknown {
		provider = config.AuthProviderMicrosoft
	}

	auth, err := h.identity.Get(provider.String())
	if err != nil {
		log.Error(err.Error(), errors.LogFields(err)...)
		return err
	}

	log.Info("performing provider logout", logger.Stringer("provider", provider))
	if err := auth.Logout(ctx); err != nil {
		log.Warn("provider logout failed", logger.Error(err))
	}

	log.Info("logout successful")
	fmt.Fprintf(opts.Stdout, "Logged out successfully\n")

	return nil
}
