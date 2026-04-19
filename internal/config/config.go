package config

// Config represents the application's top-level configuration structure.
type Config struct {
	// Auth contains settings related to authentication and provider identity.
	Auth AuthenticationConfig `json:"auth" yaml:"auth"`
	// Logging contains settings related to logging behavior and output.
	Logging LoggingConfig `json:"logging" yaml:"logging"`
	// Mounts defines the collection of virtual filesystem mount points.
	Mounts []MountConfig `json:"mounts,omitempty" yaml:"mounts,omitempty" mapstructure:"mounts"`
	// Editor contains settings related to the external editor service.
	Editor EditorConfig `json:"editor" yaml:"editor"`
}

// EditorConfig represents the configuration for the external editor.
type EditorConfig struct {
	// Command is the explicit editor command to use.
	Command string `json:"command" yaml:"command"`
}
