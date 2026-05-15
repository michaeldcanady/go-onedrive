package config

// Service coordinates the retrieval and persistence of application settings.
// It supports dot-notation for nested keys and provides fallback to system defaults.
type Service interface {
	// Get returns the value associated with the given key.
	Get(key string) (any, error)

	// Set updates the value for the given key and persists the change to the repository.
	Set(key string, value string) error

	// All returns a complete map of all user-defined configuration settings.
	All() (map[string]any, error)
}

// Repository handles the low-level persistence of configuration data.
type Repository interface {
	// Load retrieves the configuration map from the underlying storage.
	Load() (map[string]any, error)

	// Save persists the provided configuration map to the underlying storage.
	Save(config map[string]any) error
}
