package config

import (
	"context"
	"strings"

	domainconfig "github.com/michaeldcanady/go-onedrive/internal2/domain/config"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
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
		return domainconfig.ErrAlreadyRegistered
	}
	s.paths[name] = path
	return nil
}

func (s *ConfigService) getDefaultConfig() domainconfig.Configuration {
	return domainconfig.Configuration{
		Auth: domainconfig.AuthenticationConfig{
			Type:        "interactiveBrowser",
			ClientID:    "6b1e6ec0-ad93-4175-a0e0-84c02e13f206",
			TenantID:    "common",
			RedirectURI: "http://localhost:8400",
		},
	}
}

func (s *ConfigService) GetConfiguration(ctx context.Context, name string) (domainconfig.Configuration, error) {
	if err := ctx.Err(); err != nil {
		return domainconfig.Configuration{}, err
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
		return domainconfig.Configuration{}, domainconfig.ErrNotRegistered
	}

	if strings.TrimSpace(path) == "" {
		logger.Error("registered configuration path is empty",
			logging.String("event", eventConfigGetPathMissing),
		)
		return domainconfig.Configuration{}, domainconfig.ErrPathMissing
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
		return domainconfig.Configuration{}, err
	}

	logger.Info("configuration loaded successfully",
		logging.String("event", eventConfigGetLoadSuccess),
	)

	return loadedCfg, nil
}
