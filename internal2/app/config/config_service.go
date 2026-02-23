package config

import (
	"context"
	"strings"

	domainconfig "github.com/michaeldcanady/go-onedrive/internal2/domain/config"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type ConfigService struct {
	paths  map[string]string
	loader domainconfig.Loader
	logger logging.Logger
}

func New2(loader domainconfig.Loader, logger logging.Logger) *ConfigService {
	return &ConfigService{
		paths:  make(map[string]string),
		loader: loader,
		logger: logger,
	}
}

// ───────────────────────────────────────────────────────────────────────────────
// Event Taxonomy (config.service)
// ───────────────────────────────────────────────────────────────────────────────

const (
	eventConfigGetStart         = "config.get.start"
	eventConfigGetLoadStart     = "config.get.load.start"
	eventConfigGetLoadSuccess   = "config.get.load.success"
	eventConfigGetLoadFailure   = "config.get.load.failure"
	eventConfigGetNotRegistered = "config.get.not_registered"
	eventConfigGetPathMissing   = "config.get.path_missing"
)

func (s *ConfigService) AddPath(name, path string) error {
	if _, exists := s.paths[name]; exists {
		return ErrAlreadyRegistered
	}
	s.paths[name] = path
	return nil
}

func (s *ConfigService) getDefaultConfig() config.Configuration3 {
	return config.Configuration3{
		Auth: config.AuthenticationConfigImpl{
			Type:        "interactiveBrowser",
			ClientID:    "6b1e6ec0-ad93-4175-a0e0-84c02e13f206",
			TenantID:    "common",
			RedirectURI: "http://localhost:8400",
		},
	}
}

func (s *ConfigService) GetConfiguration(ctx context.Context, name string) (config.Configuration3, error) {
	if err := ctx.Err(); err != nil {
		return config.Configuration3{}, err
	}

	correlationID := util.CorrelationIDFromContext(ctx)

	logger := s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("config_name", name),
	)

	logger.Info("starting configuration retrieval",
		logging.String("event", eventConfigGetStart),
	)

	path, ok := s.paths[name]
	if !ok {
		logger.Error("configuration name not registered",
			logging.String("event", eventConfigGetNotRegistered),
		)
		return config.Configuration3{}, ErrNotRegistered
	}

	if strings.TrimSpace(path) == "" {
		logger.Error("registered configuration path is empty",
			logging.String("event", eventConfigGetPathMissing),
		)
		return config.Configuration3{}, ErrPathMissing
	}

	logger.Info("loading configuration from disk",
		logging.String("event", eventConfigGetLoadStart),
		logging.String("path", path),
	)

	loadedCfg, err := s.loader.Load(path)
	if err != nil {
		logger.Error("failed to load configuration from disk",
			logging.String("event", eventConfigGetLoadFailure),
			logging.Error(err),
		)
		return config.Configuration3{}, err
	}

	logger.Info("configuration loaded successfully",
		logging.String("event", eventConfigGetLoadSuccess),
	)

	return loadedCfg, nil
}
