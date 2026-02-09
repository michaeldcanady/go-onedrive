package config

import (
	"context"
	"errors"
	"os"
	"strings"

	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	domainconfig "github.com/michaeldcanady/go-onedrive/internal2/domain/config"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/core"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

// ConfigService manages named configuration sources and provides lazy,
// cached access to their contents.
//
// A configuration is identified by a string key (e.g. "personal")
// and associated with a filesystem path registered via AddPath.
// When GetConfiguration is called, the service attempts to retrieve
// the configuration from the CacheService. If it is not present,
// the service loads it from disk using the Loader and stores it
// in the cache for future retrieval.
//
// ConfigService does not interpret or validate configuration contents.
// It simply loads and returns them as config.Configuration3 values.
type ConfigService struct {
	cacheService domaincache.CacheService
	paths        map[string]string
	loader       domainconfig.Loader
	logger       logging.Logger
}

func New2(cache domaincache.CacheService, loader domainconfig.Loader, logger logging.Logger) *ConfigService {
	return &ConfigService{
		cacheService: cache,
		paths:        make(map[string]string),
		loader:       loader,
		logger:       logger,
	}
}

// AddPath registers a configuration source under the given name.
//
// The name is an arbitrary identifier (e.g. "default") used
// later when retrieving the configuration. The path must point to
// a readable configuration file compatible with the Loader.
//
// If a configuration with the same name is already registered,
// AddPath returns ErrAlreadyRegistered.
func (s *ConfigService) AddPath(name, path string) error {
	s.logger.Debug("registering configuration path",
		logging.String("event", "config_add_path"),
		logging.String("name", name),
		logging.String("path", path),
	)

	if _, exists := s.paths[name]; exists {
		s.logger.Warn("configuration path already registered",
			logging.String("event", "config_add_path"),
			logging.String("name", name),
		)
		return ErrAlreadyRegistered
	}

	s.paths[name] = path

	s.logger.Info("configuration path registered",
		logging.String("event", "config_add_path"),
		logging.String("name", name),
		logging.String("path", path),
	)

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

// GetConfiguration returns the configuration associated with the given name.
//
// The method first checks the provided context for cancellation.
// It then attempts to retrieve the configuration from the CacheService.
// If the configuration is not cached, the service loads it from the
// registered path using the Loader and stores it in the domaincache.
//
// Errors:
//   - ErrNotRegistered: no path has been registered for this name
//   - ErrPathMissing: the registered path is empty or whitespace
//   - Loader errors: if the configuration file cannot be read or parsed
//   - CacheService errors: if caching the loaded configuration fails
//
// The returned config.Configuration3 is the raw configuration data
// as loaded by the Loader.
func (s *ConfigService) GetConfiguration(ctx context.Context, name string) (config.Configuration3, error) {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Debug("retrieving configuration",
		logging.String("event", "config_get"),
		logging.String("name", name),
		logging.String("correlation_id", cid),
	)

	// Context canceled
	if err := ctx.Err(); err != nil {
		s.logger.Warn("context canceled while retrieving configuration",
			logging.String("event", "config_get"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return config.Configuration3{}, err
	}

	// Try cache first
	cfg, err := s.cacheService.GetConfiguration(ctx, name)
	if err == nil {
		s.logger.Debug("configuration retrieved from cache",
			logging.String("event", "config_get_cache_hit"),
			logging.String("name", name),
			logging.String("correlation_id", cid),
		)
		return cfg, nil
	}

	if errors.Is(err, domaincache.ErrUnavailableCache) {
		s.logger.Warn("cache service unavailable while retrieving configuration",
			logging.String("event", "config_get_cache_unavailable"),
			logging.Error(err),
			logging.String("name", name),
			logging.String("correlation_id", cid),
		)
	}

	// Unexpected cache error
	if !errors.Is(err, core.ErrKeyNotFound) && !errors.Is(err, domaincache.ErrUnavailableCache) {
		s.logger.Error("failed to retrieve configuration from cache",
			logging.String("event", "config_get_cache_error"),
			logging.Error(err),
			logging.String("name", name),
			logging.String("correlation_id", cid),
		)
		return config.Configuration3{}, err
	}

	// Cache miss
	s.logger.Info("configuration not found in cache",
		logging.String("event", "config_get_cache_miss"),
		logging.String("name", name),
		logging.String("correlation_id", cid),
	)

	// Validate registration
	path, ok := s.paths[name]
	if !ok {
		s.logger.Warn("configuration name not registered",
			logging.String("event", "config_get_not_registered"),
			logging.String("name", name),
			logging.String("correlation_id", cid),
		)
		return config.Configuration3{}, ErrNotRegistered
	}

	if strings.TrimSpace(path) == "" {
		s.logger.Error("registered configuration path is empty",
			logging.String("event", "config_get_invalid_path"),
			logging.String("name", name),
			logging.String("correlation_id", cid),
		)
		return config.Configuration3{}, ErrPathMissing
	}

	// Load from disk
	s.logger.Debug("loading configuration from disk",
		logging.String("event", "config_load_disk"),
		logging.String("name", name),
		logging.String("path", path),
		logging.String("correlation_id", cid),
	)

	loadedCfg, err := s.loader.Load(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			s.logger.Error("failed to load configuration from disk",
				logging.String("event", "config_load_disk"),
				logging.Error(err),
				logging.String("name", name),
				logging.String("path", path),
				logging.String("correlation_id", cid),
			)
			return config.Configuration3{}, err
		}
		s.logger.Warn("configuration file does not exist on disk, using default configuration",
			logging.String("event", "config_load_missing_file"),
			logging.String("name", name),
			logging.String("path", path),
			logging.String("correlation_id", cid),
		)
		loadedCfg = s.getDefaultConfig()
	}

	if loadedCfg == (config.Configuration3{}) {
		s.logger.Warn("loaded configuration is empty, using default configuration",
			logging.String("event", "config_load_empty"),
			logging.String("name", name),
			logging.String("path", path),
			logging.String("correlation_id", cid),
		)
		loadedCfg = s.getDefaultConfig()
	}

	// Save to cache
	s.logger.Debug("caching loaded configuration",
		logging.String("event", "config_cache_set"),
		logging.String("name", name),
		logging.String("correlation_id", cid),
	)

	if err := s.cacheService.SetConfiguration(ctx, name, loadedCfg); err != nil {
		if !errors.Is(err, domaincache.ErrUnavailableCache) {
			s.logger.Error("failed to cache configuration",
				logging.String("event", "config_cache_set"),
				logging.Error(err),
				logging.String("name", name),
				logging.String("correlation_id", cid),
			)
			return config.Configuration3{}, err
		}
		s.logger.Warn("cache service unavailable while caching configuration",
			logging.String("event", "config_cache_set_unavailable"),
			logging.Error(err),
			logging.String("name", name),
			logging.String("correlation_id", cid),
		)
	}

	s.logger.Info("configuration loaded successfully",
		logging.String("event", "config_get_success"),
		logging.String("name", name),
		logging.String("correlation_id", cid),
	)

	return loadedCfg, nil
}
