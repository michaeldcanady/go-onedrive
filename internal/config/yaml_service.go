package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"gopkg.in/yaml.v3"
)

// YAMLService is an implementation of the Service interface that loads configuration from YAML files.
type YAMLService struct {
	// mu protects the paths map from concurrent access.
	mu sync.RWMutex
	// paths maps profile names to their YAML configuration file paths (overrides).
	paths map[string]string
	// resolver is used to look up configuration paths if not in the override map.
	resolver PathResolver
	// log is the logger used for reporting configuration events.
	log logger.Logger
}

// NewYAMLService creates a new instance of YAMLService.
func NewYAMLService(resolver PathResolver, log logger.Logger) *YAMLService {
	return &YAMLService{
		paths:    make(map[string]string),
		resolver: resolver,
		log:      log,
	}
}

// AddPath registers a configuration file path for the given profile.
func (s *YAMLService) AddPath(profile, path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.paths[profile] = path
	return nil
}

// GetPath retrieves the registered configuration file path for the given profile.
func (s *YAMLService) GetPath(profile string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check for override first
	if path, ok := s.paths[profile]; ok && path != "" {
		return path, true
	}

	// Then check resolver if available
	if s.resolver != nil {
		if path, err := s.resolver.ResolvePath(context.Background(), profile); err == nil && path != "" {
			return path, true
		}
	}

	return "", false
}

// GetConfig reads and unmarshals the YAML configuration for the specified profile.
func (s *YAMLService) GetConfig(ctx context.Context, profile string) (Config, error) {
	path, ok := s.GetPath(profile)

	if !ok || path == "" {
		return s.defaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s.defaultConfig(), nil
		}
		return Config{}, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	cfg := s.defaultConfig()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

// SaveConfig writes the provided configuration to the YAML file associated with the profile.
func (s *YAMLService) SaveConfig(ctx context.Context, profile string, cfg Config) error {
	path, ok := s.GetPath(profile)

	if !ok || path == "" {
		return fmt.Errorf("no configuration path registered for profile: %s", profile)
	}

	// Ensure the parent directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create configuration directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	s.log.Debug("saved configuration", logger.String("profile", profile), logger.String("path", path))
	return nil
}

// defaultConfig returns the fallback configuration used when no file is found.
func (s *YAMLService) defaultConfig() Config {
	return Config{
		Auth: AuthenticationConfig{
			Provider:    "microsoft",
			ClientID:    "6b1e6ec0-ad93-4175-a0e0-84c02e13f206",
			TenantID:    "common",
			RedirectURI: "http://localhost:8400",
		},
	}
}
