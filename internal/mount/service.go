package mount

import (
	"context"
)

// MountConfig represents the configuration for a mount point.
type MountConfig struct {
	Path       string            `json:"path" yaml:"path"`
	Type       string            `json:"type" yaml:"type"`
	IdentityID string            `json:"identity_id,omitempty" yaml:"identity_id,omitempty"`
	Options    map[string]string `json:"options,omitempty" yaml:"options,omitempty"`
}

// Service provides methods for managing virtual filesystem mount points.
type Service interface {
	// ListMounts retrieves all configured mount points.
	ListMounts(ctx context.Context) ([]MountConfig, error)
	// AddMount adds or updates a mount point in the configuration.
	AddMount(ctx context.Context, m MountConfig) error
	// RemoveMount removes a mount point from the configuration.
	RemoveMount(ctx context.Context, path string) error
}
