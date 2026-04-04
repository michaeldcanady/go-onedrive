package profile

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/state"
)

// Service provides operations for managing user configuration profiles.
type Service interface {
	// Get retrieves a specific profile by name.
	Get(ctx context.Context, name string) (Profile, error)
	// List returns all registered profiles.
	List(ctx context.Context) ([]Profile, error)
	// Create initializes a new profile with the given name.
	Create(ctx context.Context, name string) (Profile, error)
	// Delete removes a profile and its associated data.
	Delete(ctx context.Context, name string) error
	// Exists checks if a profile with the given name already exists.
	Exists(ctx context.Context, name string) (bool, error)
	// Update saves the specified profile.
	Update(ctx context.Context, p Profile) error
	// GetActive retrieves the currently active profile.
	GetActive(ctx context.Context) (Profile, error)
	// SetActive marks a specific profile as the active one with the given scope.
	SetActive(ctx context.Context, name string, scope state.Scope) error
}

// PathResolver defines an interface for resolving profile-specific configuration paths.
type PathResolver interface {
	// ResolvePath returns the configuration file path for the specified profile name.
	ResolvePath(ctx context.Context, profileName string) (string, error)
}

// StateProvider defines the interface required by the profile service to manage its state.
type StateProvider interface {
	Get(key Key) (string, error)
	Set(key Key, value string, scope Scope) error
}

// Key identifies a piece of application state (e.g., active profile or drive).
type Key int

const (
	// KeyProfile represents the currently active profile.
	KeyProfile Key = iota
)

// Scope determines the persistence level of state data.
type Scope int

const (
	// ScopeGlobal state persists across sessions.
	ScopeGlobal Scope = iota
)
