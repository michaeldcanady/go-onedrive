package set

import (
	"context"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/identity/shared"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
)

// Handler orchestrates the config set operation.
type Handler struct {
	config  config.Service
	profile profile.Service
	log     logger.Logger
}

// NewHandler initializes a new instance of the set Handler.
func NewHandler(
	cfg config.Service,
	prof profile.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("config-set")
	return &Handler{
		config:  cfg,
		profile: prof,
		log:     cliLog,
	}
}

// Handle updates the configuration setting.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx)

	p, err := h.profile.GetActive(ctx)
	if err != nil {
		log.Error(err.Error(), errors.LogFields(err)...)
		return err
	}

	cfg, err := h.config.GetConfig(ctx)
	if err != nil {
		log.Error(err.Error(), errors.LogFields(err)...)
		return err
	}

	if err := h.setValueByKey(&cfg, opts.Key, opts.Value); err != nil {
		log.Error(err.Error(), errors.LogFields(err)...)
		return err
	}

	if err := h.config.SaveConfig(ctx, cfg); err != nil {
		log.Error(err.Error(), errors.LogFields(err)...)
		return err
	}

	fmt.Fprintf(opts.Stdout, "Updated %s for profile %s\n", opts.Key, p.Name)
	return nil
}

func (h *Handler) setValueByKey(cfg *config.Config, key, value string) error {
	key = strings.ToLower(key)
	switch key {
	case "auth.provider":
		cfg.Auth.Provider = config.ParseAuthProvider(value)
	case "auth.client_id":
		cfg.Auth.ClientID = value
	case "auth.tenant_id":
		cfg.Auth.TenantID = value
	case "auth.client_secret":
		cfg.Auth.ClientSecret = value
	case "auth.method":
		cfg.Auth.Method = shared.ParseAuthMethod(value)
	case "auth.redirect_uri":
		cfg.Auth.RedirectURI = value
	case "logging.level":
		cfg.Logging.Level = logger.ParseLevel(value)
	case "logging.output":
		cfg.Logging.Output = value
	case "logging.format":
		cfg.Logging.Format = value
	default:
		return errors.NewAppError(errors.CodeInvalidInput, nil, "unknown configuration key", "Use 'odc config get' to see all available keys.")
	}
	return nil
}
