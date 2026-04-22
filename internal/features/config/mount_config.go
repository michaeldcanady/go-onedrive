package config

// MountConfig defines a single virtual filesystem mount point.
type MountConfig struct {
	// Path is the virtual path where the backend is mounted (e.g., "/work").
	Path string `json:"path" yaml:"path" mapstructure:"path"`
	// Type is the backend type (e.g., "local", "onedrive", "grpc").
	Type string `json:"type" yaml:"type" mapstructure:"type"`
	// IdentityID is the identifier for the account/credentials to use.
	IdentityID string `json:"identity_id,omitempty" yaml:"identity_id,omitempty" mapstructure:"identity_id"`
	// Options contains backend-specific settings (e.g., "root" for local, "drive_id" for onedrive).
	Options map[string]string `json:"options,omitempty" yaml:"options,omitempty" mapstructure:"options"`
}
