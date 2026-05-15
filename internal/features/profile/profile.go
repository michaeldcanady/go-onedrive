package profile

// Profile represents a named collection of user settings and session state.
type Profile struct {
	Name string `json:"name"`
}

// Service coordinates the management of multiple [Profile] instances and
// tracks the active ('current') profile for the session.
type Service interface {
	// Create registers a new named profile.
	Create(name string) (*Profile, error)

	// List returns all registered profiles.
	List() ([]*Profile, error)

	// Delete removes a profile and its associated state from the repository.
	Delete(name string) error

	// GetCurrent retrieves the currently active profile.
	GetCurrent() (*Profile, error)

	// SetCurrent marks the specified profile as the active one for the session.
	SetCurrent(name string) error
}

// Repository handles the low-level persistence of [Profile] metadata and
// session-wide state (like the current profile indicator).
type Repository interface {
	// Create persists a new profile to the underlying store.
	Create(p *Profile) error

	// List retrieves all stored profiles.
	List() ([]*Profile, error)

	// Delete removes a profile from the underlying store.
	Delete(name string) error

	// GetCurrent retrieves the name of the currently active profile.
	GetCurrent() (string, error)

	// SetCurrent persists the name of the currently active profile.
	SetCurrent(name string) error
}
