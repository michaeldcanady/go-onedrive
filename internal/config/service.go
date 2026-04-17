package config

import (
	"context"
)

// Service provides methods to retrieve and manage application configuration.
type Service interface {
	// GetConfig retrieves the Configuration for the active profile.
	GetConfig(ctx context.Context) (Config, error)
	// GetPath retrieves the registered file path for the active profile.
	GetPath(ctx context.Context) (string, bool)
	// SaveConfig saves the Configuration for the active profile.
	SaveConfig(ctx context.Context, cfg Config) error

	// AddMount adds a new mount point to the configuration.
	AddMount(ctx context.Context, m MountConfig) error
	// RemoveMount removes a mount point from the configuration.
	RemoveMount(ctx context.Context, path string) error
}
