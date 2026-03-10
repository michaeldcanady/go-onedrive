package config

import (
	"context"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	domainconfig "github.com/michaeldcanady/go-onedrive/internal2/domain/config"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type ConfigService struct {
	paths  map[string]string
	loader domainconfig.Loader
	log    logger.Logger
}

func New2(loader domainconfig.Loader, log logger.Logger) *ConfigService {
	return &ConfigService{
		paths:  make(map[string]string),
		loader: loader,
		log:    log,
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
			Type:        auth.MethodInteractiveBrowser,
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

	log := s.log.WithContext(ctx).With(
		logger.String("correlation_id", correlationID),
		logger.String("config_name", name),
	)

	log.Info("starting configuration retrieval",
		logger.String("event", eventConfigGetStart),
	)

	path, ok := s.paths[name]
	if !ok {
		log.Error("configuration name not registered",
			logger.String("event", eventConfigGetNotRegistered),
		)
		return domainconfig.Configuration{}, domainconfig.ErrNotRegistered
	}

	if strings.TrimSpace(path) == "" {
		log.Error("registered configuration path is empty",
			logger.String("event", eventConfigGetPathMissing),
		)
		return domainconfig.Configuration{}, domainconfig.ErrPathMissing
	}

	log.Info("loading configuration from disk",
		logger.String("event", eventConfigGetLoadStart),
		logger.String("path", path),
	)

	loadedCfg, err := s.loader.Load(path)
	if err != nil {
		log.Error("failed to load configuration from disk",
			logger.String("event", eventConfigGetLoadFailure),
			logger.Error(err),
		)
		return domainconfig.Configuration{}, err
	}

	log.Info("configuration loaded successfully",
		logger.String("event", eventConfigGetLoadSuccess),
	)

	return loadedCfg, nil
}
