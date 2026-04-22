package profile

import "time"

// Profile represents a user configuration profile.
type Profile struct {
	// Name is the unique identifier for the profile.
	Name string `json:"name"`
	// ConfigPath is the absolute path to the configuration file for this profile.
	ConfigPath string `json:"config_path,omitempty"`
	// CreatedAt is the timestamp when the profile was first created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the timestamp of the last modification to the profile.
	UpdatedAt time.Time `json:"updated_at"`
}
