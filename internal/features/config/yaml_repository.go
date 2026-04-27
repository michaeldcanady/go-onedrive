package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"gopkg.in/yaml.v3"
)

// YAMLRepository implements the Repository interface using YAML files.
type YAMLRepository struct {
	path string
	log  logger.Logger
}

// NewYAMLRepository creates a new instance of YAMLRepository.
func NewYAMLRepository(path string, log logger.Logger) *YAMLRepository {
	return &YAMLRepository{
		path: path,
		log:  log,
	}
}

// Load reads and unmarshals the YAML configuration.
func (r *YAMLRepository) Load(ctx context.Context) (*Config, error) {
	l := r.log.WithContext(ctx)

	if r.path == "" {
		return r.defaultConfig(), nil
	}

	data, err := os.ReadFile(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return r.defaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config file %s: %w", r.path, err)
	}

	cfg := r.defaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	l.Debug("configuration loaded", logger.String("path", r.path))
	return cfg, nil
}

// Save persists the configuration to a YAML file.
func (r *YAMLRepository) Save(ctx context.Context, cfg *Config) error {
	if r.path == "" {
		return fmt.Errorf("no configuration path provided")
	}

	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(r.path, data, 0o600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// defaultConfig returns the fallback configuration.
func (r *YAMLRepository) defaultConfig() *Config {
	return &Config{
		Auth: AuthenticationConfig{
			Provider:    "microsoft",
			ClientID:    "6b1e6ec0-ad93-4175-a0e0-84c02e13f206",
			TenantID:    "common",
			RedirectURI: "http://localhost:8400",
		},
		Logging: LoggingConfig{
			Level: logger.LevelInfo,
		},
	}
}
