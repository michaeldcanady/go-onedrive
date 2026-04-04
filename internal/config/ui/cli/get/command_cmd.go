package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/state"
	"gopkg.in/yaml.v3"
)

// Handler orchestrates the config get operation.
type Handler struct {
	config config.Service
	state  state.Service
	log    logger.Logger
}

// NewHandler initializes a new instance of the get Handler.
func NewHandler(
	cfg config.Service,
	st state.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("config-get")
	return &Handler{
		config: cfg,
		state:  st,
		log:    cliLog,
	}
}

// Handle retrieves and displays the configuration.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	cfg, err := h.config.GetConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if opts.Key == "" {
		data, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("failed to marshal configuration: %w", err)
		}
		fmt.Fprint(opts.Stdout, string(data))
		return nil
	}

	value, err := h.getValueByKey(cfg, opts.Key)
	if err != nil {
		return err
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
	default:
		return nil, fmt.Errorf("unknown configuration key: %s", key)
	}
}
