package registry

import (
	"context"
	"fmt"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal/identity/shared"
)

// Registry is a thread-safe implementation of the Service interface.
type Registry struct {
	// mu protects the providers map.
	mu sync.RWMutex
	// providers maps provider names to their respective authenticators.
	providers map[string]shared.Authenticator
}

// NewRegistry initializes a new instance of the Registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]shared.Authenticator),
	}
}

// Register adds an authenticator to the registry.
func (r *Registry) Register(provider string, auth shared.Authenticator) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[provider] = auth
}

// Get retrieves the authenticator for the specified provider.
func (r *Registry) Get(provider string) (shared.Authenticator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	auth, ok := r.providers[provider]
	if !ok {
		return nil, fmt.Errorf("identity provider %s not found in registry", provider)
	}
	return auth, nil
}

// ListIdentities returns all identity IDs for a specific provider.
func (r *Registry) ListIdentities(ctx context.Context, provider string) ([]string, error) {
	auth, err := r.Get(provider)
	if err != nil {
		return nil, err
	}
	return auth.ListIdentities(ctx)
}
