package config

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"gopkg.in/yaml.v3"
)

// YAMLService is an implementation of the Service interface that loads configuration from YAML files.
type YAMLService struct {
	// mu protects the paths map from concurrent access.
	mu sync.RWMutex
	// paths maps profile names to their YAML configuration file paths.
	paths map[string]string
	// log is the logger used for reporting configuration events.
	log logger.Logger
}

// NewYAMLService creates a new instance of YAMLService.
func NewYAMLService(log logger.Logger) *YAMLService {
	return &YAMLService{
		paths: make(map[string]string),
		log:   log,
	}
}

// AddPath registers a configuration file path for the given profile.
func (s *YAMLService) AddPath(profile, path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.paths[profile] = path
	return nil
}

// GetConfig reads and unmarshals the YAML configuration for the specified profile.
func (s *YAMLService) GetConfig(ctx context.Context, profile string) (Config, error) {
	s.mu.RLock()
	path, ok := s.paths[profile]
	s.mu.RUnlock()

	if !ok || path == "" {
		return s.defaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s.defaultConfig(), nil
		}
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
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
