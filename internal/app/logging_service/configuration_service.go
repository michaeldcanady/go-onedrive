package loggingservice

import "context"

// ConfigurationService provides access to configuration values.
type ConfigurationService interface {
	// GetString retrieves a string value from the configuration by key.
	GetString(context.Context, string) (string, error)
	// GetStringDefault retrieves a string value from the configuration by key,
	GetStringDefault(context.Context, string, string) (string, error)
}
