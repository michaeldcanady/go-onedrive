package registry

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs/shared"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/core/state"
)

// Service defines the interface for managing and resolving filesystem providers.
type Service interface {
	// Register associates a service with its provider name.
	Register(provider string, svc shared.Service)
	// Get retrieves the service for a given provider name.
	Get(provider string) (shared.Service, error)
	// Resolve identifies the provider for a path based on its prefix.
	Resolve(ctx context.Context, path string) (shared.Service, string, error)
}

// Registry manages a collection of filesystem providers.
type Registry struct {
	mu sync.RWMutex
	// providers maps provider names to their implementations.
	providers map[string]shared.Service
	// state is the service for tracking the active provider and aliases.
	state state.Service
	// logger is the logger instance.
	logger logger.Logger
}

// NewRegistry initializes a new instance of the Registry.
func NewRegistry(state state.Service, log logger.Logger) *Registry {
	return &Registry{
		providers: make(map[string]shared.Service),
		state:     state,
		logger:    log,
	}
}

// Register associates a service with its provider name.
func (r *Registry) Register(provider string, svc shared.Service) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Wrap the service with the validation decorator if it's a filesystem service
	r.providers[provider] = shared.NewValidationDecorator(svc, r.logger)
}

// Get retrieves the service for a given provider name.
func (r *Registry) Get(provider string) (shared.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	svc, ok := r.providers[provider]
	if !ok {
		return nil, fmt.Errorf("filesystem provider %s not found in registry", provider)
	}
	return svc, nil
}

// Resolve identifies the provider for a path based on its prefix (e.g., "onedrive:path").
// If no prefix is present, it defaults to "onedrive".
func (r *Registry) Resolve(ctx context.Context, path string) (shared.Service, string, error) {
	prefix, rest, found := strings.Cut(path, ":")
	if !found {
		// No prefix, default to onedrive
		p, err := r.Get(DefaultProviderPrefix)
		if err != nil {
			return nil, "", err
		}
		return p, path, nil
	}

	// Check if prefix is a registered provider
	p, err := r.Get(prefix)
	if err != nil {
		// Check if prefix is an alias
		driveID, err := r.state.GetDriveAlias(prefix)
		if err != nil {
			return nil, "", fmt.Errorf("unknown provider or alias: %s", prefix)
		}
		// If it's an alias, use the default provider (onedrive) and prepend the drive ID
		defaultProvider, err := r.Get(DefaultProviderPrefix)
		if err != nil {
			return nil, "", err
		}
		// Construct path with drive ID for the default provider
		rest = fmt.Sprintf("%s:%s", driveID, rest)
		return defaultProvider, rest, nil
	}
	// If it's a registered provider, use it directly
	return p, rest, nil
}

// DefaultProviderPrefix is the default provider to use when no prefix is specified.
const DefaultProviderPrefix = "onedrive"

// Cleanup removes all providers from the registry.
func (r *Registry) Cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers = make(map[string]shared.Service)
}
