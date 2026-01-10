package di

import "context"

// ConfigurationService provides access to configuration values.
type ConfigurationService interface {
	// GetString retrieves a string value from the configuration by key.
	GetString(context.Context, string) (string, error)
	// GetStringDefault retrieves a string value from the configuration by key,
	GetStringDefault(context.Context, string, string) (string, error)

	// SetConfigFile sets the configuration file path.
	SetConfigFile(ctx context.Context, path string)

	// WriteConfiguration writes the current configuration to file.
	WriteConfiguration(ctx context.Context) error
}
