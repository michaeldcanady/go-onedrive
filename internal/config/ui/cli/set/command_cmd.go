package set

import (
	"context"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/state"
)

// Handler orchestrates the config set operation.
type Handler struct {
	config config.Service
	state  state.Service
	log    logger.Logger
}

// NewHandler initializes a new instance of the set Handler.
func NewHandler(
	cfg config.Service,
	st state.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("config-set")
	return &Handler{
		config: cfg,
		state:  st,
		log:    cliLog,
	}
}

// Handle updates the configuration setting.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	profileName, err := h.state.Get(state.KeyProfile)
	if err != nil {
		return fmt.Errorf("failed to get current profile: %w", err)
	}

	cfg, err := h.config.GetConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configuration for profile %s: %w", profileName, err)
	}

	if err := h.setValueByKey(&cfg, opts.Key, opts.Value); err != nil {
		return err
	}

	if err := h.config.SaveConfig(ctx, cfg); err != nil {
		return fmt.Errorf("failed to save configuration for profile %s: %w", profileName, err)
	}

	fmt.Fprintf(opts.Stdout, "Updated %s for profile %s\n", opts.Key, profileName)
	return nil
}

func (h *Handler) setValueByKey(cfg *config.Config, key, value string) error {
	key = strings.ToLower(key)
	switch key {
	case "auth.provider":
		cfg.Auth.Provider = value
	case "auth.client_id":
		cfg.Auth.ClientID = value
	case "auth.tenant_id":
		cfg.Auth.TenantID = value
	case "auth.client_secret":
		cfg.Auth.ClientSecret = value
	case "auth.method":
		cfg.Auth.Method = value
	case "auth.redirect_uri":
		cfg.Auth.RedirectURI = value
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}
	return nil
}
