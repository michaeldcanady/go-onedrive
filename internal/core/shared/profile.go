package shared

// Profile represents a user's configuration profile.
// It contains metadata such as the profile name and associated file paths.
type Profile struct {
	// Name is the unique identifier for the profile.
	Name string `json:"name" yaml:"name"`
	// Path is the filesystem path where profile data is stored.
	Path string `json:"path" yaml:"path"`
	// ConfigPath is the location of the profile's specific configuration file.
	ConfigPath string `json:"config_path" yaml:"config_path"`
}
