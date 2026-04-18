package mount

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/config"
)

// Service provides methods for managing virtual filesystem mount points.
type Service interface {
	// ListMounts retrieves all configured mount points.
	ListMounts(ctx context.Context) ([]config.MountConfig, error)
	// AddMount adds or updates a mount point in the configuration.
	AddMount(ctx context.Context, m config.MountConfig) error
	// RemoveMount removes a mount point from the configuration.
	RemoveMount(ctx context.Context, path string) error
}
