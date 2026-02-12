package config

import (
	"context"
	"errors"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	domainconfig "github.com/michaeldcanady/go-onedrive/internal2/domain/config"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/core"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type ConfigService struct {
	cacheService cache.CacheService
	paths        map[string]string
	loader       domainconfig.Loader
	logger       logging.Logger
}

func New2(cache cache.CacheService, loader domainconfig.Loader, logger logging.Logger) *ConfigService {
	return &ConfigService{
		cacheService: cache,
		paths:        make(map[string]string),
		loader:       loader,
		logger:       logger,
	}
}

// ───────────────────────────────────────────────────────────────────────────────
// Event Taxonomy (config.service)
// ───────────────────────────────────────────────────────────────────────────────

const (
	eventConfigGetStart         = "config.get.start"
	eventConfigGetCacheHit      = "config.get.cache.hit"
	eventConfigGetCacheMiss     = "config.get.cache.miss"
	eventConfigGetCacheEmpty    = "config.get.cache.empty"
	eventConfigGetLoadStart     = "config.get.load.start"
	eventConfigGetLoadSuccess   = "config.get.load.success"
	eventConfigGetLoadFailure   = "config.get.load.failure"
	eventConfigGetSaveStart     = "config.get.save.start"
	eventConfigGetSaveSuccess   = "config.get.save.success"
	eventConfigGetSaveFailure   = "config.get.save.failure"
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

	cfg, err := s.cacheService.GetConfiguration(ctx, name)
	if err == nil {
		logger.Debug("configuration retrieved from cache",
			logging.String("event", eventConfigGetCacheHit),
		)

		if cfg == (config.Configuration3{}) {
			logger.Warn("cached configuration is empty; using default",
				logging.String("event", eventConfigGetCacheEmpty),
			)
			return s.getDefaultConfig(), nil
		}

		return cfg, nil
	}

	if !errors.Is(err, core.ErrKeyNotFound) {
		logger.Error("failed to retrieve configuration from cache",
			logging.String("event", eventConfigGetCacheMiss),
			logging.Error(err),
		)
		return config.Configuration3{}, err
	}

	logger.Info("configuration not found in cache",
		logging.String("event", eventConfigGetCacheMiss),
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

	logger.Debug("saving configuration to cache",
		logging.String("event", eventConfigGetSaveStart),
	)

	if err := s.cacheService.SetConfiguration(ctx, name, loadedCfg); err != nil {
		logger.Error("failed to save configuration to cache",
			logging.String("event", eventConfigGetSaveFailure),
			logging.Error(err),
		)
		return config.Configuration3{}, err
	}

	logger.Info("configuration cached successfully",
		logging.String("event", eventConfigGetSaveSuccess),
	)

	return loadedCfg, nil
}
