package state

import "context"

// Service provides methods to retrieve, update, and manage the application state.
// It is a composite interface of all domain-specific stores.
type Service interface {
	ProfileStore
	DriveStore
	AuthStore
	ConfigStore

	// Get retrieves a raw state value by its key.
	Get(key Key) (string, error)
	// Set assigns a raw value to a key within the specified scope.
	Set(key Key, value string, scope Scope) error
	// Clear removes a state value for the given key from all scopes.
	Clear(key Key) error
}

// ProfileStore manages the active profile state.
type ProfileStore interface {
	// GetProfile retrieves the name of the currently active profile.
	GetProfile(ctx context.Context) (string, error)
	// SetProfile updates the name of the active profile.
	SetProfile(ctx context.Context, name string, scope Scope) error
}

// DriveStore manages the active drive state.
type DriveStore interface {
	// GetDrive retrieves the ID of the currently active drive.
	GetDrive(ctx context.Context) (string, error)
	// SetDrive updates the ID of the active drive.
	SetDrive(ctx context.Context, driveID string, scope Scope) error
}

// AuthStore manages authentication-related state, such as access tokens.
type AuthStore interface {
	// GetAccessToken retrieves the cached authentication token.
	GetAccessToken(ctx context.Context) (string, error)
	// SetAccessToken updates the cached authentication token.
	SetAccessToken(ctx context.Context, token string, scope Scope) error
	// ClearAccessToken removes the cached authentication token.
	ClearAccessToken(ctx context.Context) error
}

// ConfigStore manages transient configuration overrides.
type ConfigStore interface {
	// GetConfigOverride retrieves the path to a configuration file override.
	GetConfigOverride(ctx context.Context) (string, error)
	// SetConfigOverride updates the configuration file path override.
	SetConfigOverride(ctx context.Context, path string, scope Scope) error
}
