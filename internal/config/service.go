package config

import (
	"context"
)

// Service provides methods to retrieve and manage application configuration.
type Service interface {
	// GetConfig retrieves the Configuration for a specific profile.
	GetConfig(ctx context.Context, profile string) (Config, error)
	// AddPath registers a file path for a configuration profile.
	AddPath(profile, path string) error
	// GetPath retrieves the registered file path for a configuration profile.
	GetPath(profile string) (string, bool)
	// SaveConfig saves the Configuration for a specific profile.
	SaveConfig(ctx context.Context, profile string, cfg Config) error
}

// PathResolver defines an interface for resolving profile-specific configuration paths.
type PathResolver interface {
	// ResolvePath returns the configuration file path for the specified profile name.
	ResolvePath(ctx context.Context, profileName string) (string, error)
}
