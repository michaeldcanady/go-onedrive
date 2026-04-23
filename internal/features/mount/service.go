package mount

import (
	"context"
)

// Service provides methods for managing virtual filesystem mount points.
type Service interface {
	// ListMounts retrieves all configured mount points.
	ListMounts(ctx context.Context) ([]MountConfig, error)
	// AddMount adds or updates a mount point in the configuration.
	AddMount(ctx context.Context, m MountConfig) error
	// RemoveMount removes a mount point from the configuration.
	RemoveMount(ctx context.Context, path string) error
	// RegisterValidator registers a validator for a given mount type.
	RegisterValidator(mountType string, v OptionValidator)
	// RegisterCompletionProvider registers a completion provider for a given mount type.
	RegisterCompletionProvider(mountType string, p CompletionProvider)
	// GetCompletionProvider retrieves a registered completion provider.
	GetCompletionProvider(mountType string) (CompletionProvider, bool)
	// GetMountOptions retrieves all registered mount options.
	GetMountOptions() map[string][]MountOption
}
