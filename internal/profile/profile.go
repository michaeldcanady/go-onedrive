package profile

import "time"

// TODO: consolidate DefaultProfileName somewhere
const (
	// DefaultProfileName is the name of the fallback profile.
	DefaultProfileName = "default"
)

// Profile represents a user configuration profile.
type Profile struct {
	// Name is the unique identifier for the profile.
	Name string `json:"name"`
	// ConfigPath is the absolute path to the configuration file for this profile.
	ConfigPath string `json:"config_path,omitempty"`
	// ActiveDriveID is the ID of the drive currently used by this profile.
	ActiveDriveID string `json:"active_drive_id,omitempty"`
	// CreatedAt is the timestamp when the profile was first created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the timestamp of the last modification to the profile.
	UpdatedAt time.Time `json:"updated_at"`
}
