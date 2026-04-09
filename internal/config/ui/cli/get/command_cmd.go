package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"gopkg.in/yaml.v3"
)

// Handler orchestrates the config get operation.
type Handler struct {
	config config.Service
	log    logger.Logger
}

// NewHandler initializes a new instance of the get Handler.
func NewHandler(
	cfg config.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("config-get")
	return &Handler{
		config: cfg,
		log:    cliLog,
	}
}

// Handle retrieves and displays the configuration.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx)

	cfg, err := h.config.GetConfig(ctx)
	if err != nil {
		log.Error(err.Error(), errors.LogFields(err)...)
		return err
	}

	if opts.Key == "" {
		data, err := yaml.Marshal(cfg)
		if err != nil {
			appErr := errors.NewAppError(errors.CodeInternal, err, "failed to marshal configuration", "")
			log.Error(appErr.Error(), errors.LogFields(appErr)...)
			return appErr
		}
		fmt.Fprint(opts.Stdout, string(data))
		return nil
	}

	value, err := h.getValueByKey(cfg, opts.Key)
	if err != nil {
		appErr := errors.NewAppError(errors.CodeInvalidInput, err, "unknown configuration key", "Use 'odc config get' without a key to see all available keys.")
		log.Error(appErr.Error(), errors.LogFields(appErr)...)
		return appErr
	}

	fmt.Fprintf(opts.Stdout, "%s: %v\n", opts.Key, value)
	return nil
}

func (h *Handler) getValueByKey(cfg config.Config, key string) (interface{}, error) {
	key = strings.ToLower(key)
	switch key {
	case "auth.provider":
		return cfg.Auth.Provider, nil
	case "auth.client_id":
		return cfg.Auth.ClientID, nil
	case "auth.tenant_id":
		return cfg.Auth.TenantID, nil
	case "auth.client_secret":
		return cfg.Auth.ClientSecret, nil
	case "auth.method":
		return cfg.Auth.Method, nil
	case "auth.redirect_uri":
		return cfg.Auth.RedirectURI, nil
	case "logging.level":
		return cfg.Logging.Level.String(), nil
	case "logging.output":
		return cfg.Logging.Output, nil
	case "logging.format":
		return cfg.Logging.Format, nil
	default:
		return nil, fmt.Errorf("unknown configuration key: %s", key)
	}
}
