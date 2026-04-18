package profile

import (
	"context"
)

// Repository defines the persistence interface for user profiles.
type Repository interface {
	// Get retrieves a profile by name.
	Get(ctx context.Context, name string) (Profile, error)
	// Create persists a new profile.
	Create(ctx context.Context, p Profile) error
	// Update modifies an existing profile.
	Update(ctx context.Context, p Profile) error
	// Delete removes a profile by name.
	Delete(ctx context.Context, name string) error
	// List retrieves all profiles.
	List(ctx context.Context) ([]Profile, error)
	// Exists checks if a profile exists by name.
	Exists(ctx context.Context, name string) (bool, error)
}
