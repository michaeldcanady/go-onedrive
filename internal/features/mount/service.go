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

type MountOption struct {
	Key    string
	Values []string
}

// OptionValidator defines the interface for backends that support option validation.
type OptionValidator interface {
	ValidateOptions(opts map[string]string) error
}

// OptionsProvider defines the interface for backends that support options
type OptionsProvider interface {
	ProvideOptions() []MountOption
}

// CompletionProvider defines the interface for backends that support dynamic completion.
type CompletionProvider interface {
	GetOptionCompletions(ctx context.Context, identityID string, toComplete string) ([]string, error)
}

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
