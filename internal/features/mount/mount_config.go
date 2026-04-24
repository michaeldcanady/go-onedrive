package mount

// MountConfig represents the configuration for a mount point.
type MountConfig struct {
	Path       string            `json:"path" yaml:"path"`
	Type       string            `json:"type" yaml:"type"`
	IdentityID string            `json:"identity_id,omitempty" yaml:"identity_id,omitempty"`
	Options    map[string]string `json:"options,omitempty" yaml:"options,omitempty"`
}
