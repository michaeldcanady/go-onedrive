package configurationservice

import (
	"context"
	"errors"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/cachev2/core"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
)

// Service2 manages named configuration sources and provides lazy,
// cached access to their contents.
//
// A configuration is identified by a string key (e.g. "personal")
// and associated with a filesystem path registered via AddPath.
// When GetConfiguration is called, the service attempts to retrieve
// the configuration from the CacheService. If it is not present,
// the service loads it from disk using the Loader and stores it
// in the cache for future retrieval.
//
// Service2 does not interpret or validate configuration contents.
// It simply loads and returns them as config.Configuration3 values.
type Service2 struct {
	cacheService CacheService
	paths        map[string]string
	loader       Loader
	logger       logging.Logger
}

func New2(cache CacheService, loader Loader, logger logging.Logger) *Service2 {
	return &Service2{
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
func (s *Service2) AddPath(name, path string) error {
	if _, exists := s.paths[name]; exists {
		return ErrAlreadyRegistered
	}
	s.paths[name] = path
	return nil
}

// GetConfiguration returns the configuration associated with the given name.
//
// The method first checks the provided context for cancellation.
// It then attempts to retrieve the configuration from the CacheService.
// If the configuration is not cached, the service loads it from the
// registered path using the Loader and stores it in the cache.
//
// Errors:
//   - ErrNotRegistered: no path has been registered for this name
//   - ErrPathMissing: the registered path is empty or whitespace
//   - Loader errors: if the configuration file cannot be read or parsed
//   - CacheService errors: if caching the loaded configuration fails
//
// The returned config.Configuration3 is the raw configuration data
// as loaded by the Loader.
func (s *Service2) GetConfiguration(ctx context.Context, name string) (config.Configuration3, error) {
	if err := ctx.Err(); err != nil {
		return config.Configuration3{}, err
	}

	cfg, err := s.cacheService.GetConfiguration(ctx, name)
	if err == nil {
		return cfg, nil
	}

	if !errors.Is(err, core.ErrKeyNotFound) {
		return config.Configuration3{}, err
	}

	s.logger.Info("config not already cached", logging.String("name", name))
	path, ok := s.paths[name]
	if !ok {
		return config.Configuration3{}, ErrNotRegistered
	}
	if strings.TrimSpace(path) == "" {
		return config.Configuration3{}, ErrPathMissing
	}

	loadedCfg, err := s.loader.Load(path)
	if err != nil {
		return config.Configuration3{}, err
	}

	if err := s.cacheService.SetConfiguration(ctx, name, loadedCfg); err != nil {
		return config.Configuration3{}, err
	}

	return loadedCfg, nil
}
