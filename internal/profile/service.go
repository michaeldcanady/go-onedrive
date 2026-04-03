package profile

import (
	"context"
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
}
