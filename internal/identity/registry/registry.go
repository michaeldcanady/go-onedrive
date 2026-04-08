package registry

import (
	"fmt"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal/identity/shared"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Registry is a thread-safe implementation of the Service interface.
type Registry struct {
	// mu protects the providers map.
	mu sync.RWMutex
	// providers maps provider names to their respective authenticators.
	providers map[string]shared.Authenticator
	log       logger.Logger
}

// NewRegistry initializes a new instance of the Registry.
func NewRegistry(l logger.Logger) *Registry {
	return &Registry{
		providers: make(map[string]shared.Authenticator),
		log:       l,
	}
}

// Register adds an authenticator to the registry.
func (r *Registry) Register(provider string, auth shared.Authenticator) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.log.Debug("registering identity provider", logger.String("provider", provider))
	r.providers[provider] = auth
}

// Get retrieves the authenticator for the specified provider.
func (r *Registry) Get(provider string) (shared.Authenticator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	auth, ok := r.providers[provider]
	if !ok {
		r.log.Warn("identity provider not found", logger.String("provider", provider))
		return nil, fmt.Errorf("identity provider %s not found in registry", provider)
	}
	r.log.Debug("retrieved identity provider", logger.String("provider", provider))
	return auth, nil
}
