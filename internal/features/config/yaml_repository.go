package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type yamlRepository struct {
	path string
}

// NewYAMLRepository creates a new YAML-based configuration repository.
func NewYAMLRepository(path string) Repository {
	return &yamlRepository{path: path}
}

func (r *yamlRepository) Load() (map[string]any, error) {
	if _, err := os.Stat(r.path); os.IsNotExist(err) {
		return make(map[string]any), nil
	}

	data, err := os.ReadFile(r.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config map[string]any
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

func (r *yamlRepository) Save(config map[string]any) error {
	dir := filepath.Dir(r.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(r.path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
