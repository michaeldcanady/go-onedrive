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
}
