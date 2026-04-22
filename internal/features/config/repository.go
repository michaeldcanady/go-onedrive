package config

import (
	"context"
)

// Repository defines the interface for persistence operations on configuration data.
type Repository interface {
	// Load retrieves the current configuration from the persistence layer.
	Load(ctx context.Context) (*Config, error)
	// Save persists the provided configuration.
	Save(ctx context.Context, cfg *Config) error
}
