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
}

// PathResolver defines an interface for resolving profile-specific configuration paths.
type PathResolver interface {
	// ResolvePath returns the configuration file path for the specified profile name.
	ResolvePath(ctx context.Context, profileName string) (string, error)
}
