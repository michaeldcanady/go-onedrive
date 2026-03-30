package registry

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs/shared"
	"github.com/michaeldcanady/go-onedrive/internal/core/state"
)

const (
	defaultProviderPrefix = "onedrive"
)

// Service manages a collection of filesystem providers, resolving them by name.
type Service interface {
	// Register associates a name with a filesystem provider.
	Register(name string, provider shared.Service)
	// Get retrieves a provider by its name.
	Get(name string) (shared.Service, error)
	// Resolve returns the appropriate provider for a given path.
	Resolve(ctx context.Context, path string) (shared.Service, string, error)
}

// Registry is a concrete implementation of the filesystem registry service.
type Registry struct {
	mu        sync.RWMutex
	providers map[string]shared.Service
	state     state.Service
}

// NewRegistry initializes a new instance of the Registry.
func NewRegistry(state state.Service) *Registry {
	return &Registry{
		providers: make(map[string]shared.Service),
		state:     state,
	}
}

// Register adds a provider to the registry.
func (r *Registry) Register(name string, provider shared.Service) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[strings.ToLower(name)] = provider
}

// Get returns a provider by its name.
func (r *Registry) Get(name string) (shared.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("filesystem provider %s not found", name)
	}
	return p, nil
}

// Resolve identifies the provider for a path based on its prefix (e.g., "onedrive:path").
// If no prefix is present, it defaults to the active drive if set, otherwise "onedrive".
func (r *Registry) Resolve(ctx context.Context, path string) (shared.Service, string, error) {
	prefix, rest, found := strings.Cut(path, ":")
	if !found {
		// No prefix, use active drive or default to onedrive
		driveID, err := r.state.Get(state.KeyDrive)
		if err == nil && driveID != "" {
			// If we have an active drive, we need to decide if we should use it.
			// For now, if no prefix, we assume it's relative to the active drive in OneDrive.
			p, err := r.Get(defaultProviderPrefix)
			return p, path, err
		}

		p, err := r.Get(defaultProviderPrefix)
		return p, path, err
	}

	// Check if prefix is a registered provider
	if p, err := r.Get(prefix); err == nil {
		rest = strings.TrimPrefix(rest, "//")
		return p, rest, nil
	}

	// Check if prefix is an alias
	if _, err := r.state.GetDriveAlias(prefix); err == nil {
		p, err := r.Get(defaultProviderPrefix)
		if err != nil {
			return nil, "", err
		}
		// Alias points to a OneDrive drive ID
		// The path should be interpreted as being within that drive.
		// For now, our providers don't yet support specifying a drive ID in the path easily,
		// but this is where that logic would go.
		rest = strings.TrimPrefix(rest, "//")
		return p, rest, nil
	}

	return nil, "", fmt.Errorf("unknown provider or alias: %s", prefix)
}
