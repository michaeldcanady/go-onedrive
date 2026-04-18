package config

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
	"github.com/michaeldcanady/go-onedrive/internal/state"
)

// ConfigService is an implementation of the Service interface that uses a Repository for persistence.
type ConfigService struct {
	repo     Repository
	resolver profile.PathResolver
	state    state.Service
	log      logger.Logger
}

// NewConfigService creates a new instance of ConfigService.
func NewConfigService(resolver profile.PathResolver, state state.Service, log logger.Logger) *ConfigService {
	return &ConfigService{
		resolver: resolver,
		state:    state,
		log:      log,
	}
}

// getRepo returns the repository initialized with the currently resolved path.
func (s *ConfigService) getRepo(ctx context.Context) (Repository, error) {
	path, _ := s.GetPath(ctx)
	return NewYAMLRepository(path, s.log), nil
}

// GetConfig retrieves the Configuration.
func (s *ConfigService) GetConfig(ctx context.Context) (Config, error) {
	repo, err := s.getRepo(ctx)
	if err != nil {
		return Config{}, err
	}
	cfg, err := repo.Load(ctx)
	if err != nil {
		return Config{}, err
	}
	return *cfg, nil
}

// GetPath retrieves the registered file path.
func (s *ConfigService) GetPath(ctx context.Context) (string, bool) {
	if path, err := s.state.Get(state.KeyConfigOverride); err == nil && path != "" {
		return path, true
	}

	profileName, err := s.state.Get(state.KeyProfile)
	if err != nil {
		return "", false
	}

	if s.resolver != nil {
		path, err := s.resolver.ResolvePath(ctx, profileName)
		if err == nil && path != "" {
			return path, true
		}
	}

	return "", false
}

// SaveConfig saves the Configuration.
func (s *ConfigService) SaveConfig(ctx context.Context, cfg Config) error {
	repo, err := s.getRepo(ctx)
	if err != nil {
		return err
	}
	return repo.Save(ctx, &cfg)
}


