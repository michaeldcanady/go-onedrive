package mount

import "context"

// ConfigRepository defines the interface for accessing mount configuration.
type ConfigRepository interface {
	GetMounts(ctx context.Context) ([]MountConfig, error)
	SaveMounts(ctx context.Context, mounts []MountConfig) error
}
