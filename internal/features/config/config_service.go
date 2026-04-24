package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/profile"
)

type ProfileProvider interface {
	GetActive(ctx context.Context) (profile.Profile, error)
}

// ConfigService is an implementation of the Service interface that uses a Repository for persistence.
type ConfigService struct {
	profileSvc ProfileProvider
	log        logger.Logger

	mu             sync.RWMutex
	configOverride string
}

// NewConfigService creates a new instance of ConfigService.
func NewConfigService(profileSvc ProfileProvider, log logger.Logger) *ConfigService {
	return &ConfigService{
		profileSvc: profileSvc,
		log:        log,
	}
}

// getRepo returns the repository initialized with the currently resolved path.
func (s *ConfigService) getRepo(ctx context.Context) (Repository, error) {
	path, ok := s.GetPath(ctx)
	if !ok {
		return nil, fmt.Errorf("could not resolve configuration path")
	}
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
	s.mu.RLock()
	override := s.configOverride
	s.mu.RUnlock()

	if override != "" {
		return override, true
	}

	if s.profileSvc != nil {
		p, err := s.profileSvc.GetActive(ctx)
		if err == nil && p.ConfigPath != "" {
			return p.ConfigPath, true
		}
	}

	return "", false
}

// UpdateConfig updates a configuration setting.
func (s *ConfigService) UpdateConfig(ctx context.Context, key string, value string) error {
	cfg, err := s.GetConfig(ctx)
	if err != nil {
		return err
	}

	if err := cfg.SetValue(key, value); err != nil {
		return fmt.Errorf("configuration update failed: %w", err)
	}

	return s.SaveConfig(ctx, cfg)
}

// SaveConfig saves the Configuration.
func (s *ConfigService) SaveConfig(ctx context.Context, cfg Config) error {
	repo, err := s.getRepo(ctx)
	if err != nil {
		return err
	}
	return repo.Save(ctx, &cfg)
}

// SetOverride sets a transient configuration path override.
func (s *ConfigService) SetOverride(ctx context.Context, path string) error {
	s.mu.Lock()
	s.configOverride = path
	s.mu.Unlock()
	return nil
}
