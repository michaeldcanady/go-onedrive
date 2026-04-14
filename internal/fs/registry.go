package fs

import (
	"fmt"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal/drive/alias"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/state"
)

// Registry manages a collection of filesystem providers.
type Registry struct {
	mu sync.RWMutex
	// providers maps provider names to their implementations.
	providers map[string]Service
	// state is the service for tracking the active provider and aliases.
	state state.Service
	// alias is the drive alias management service.
	alias alias.Service
	// logger is the logger instance.
	logger logger.Logger
}

// NewRegistry initializes a new instance of the Registry.
func NewRegistry(state state.Service, alias alias.Service, log logger.Logger) *Registry {
	return &Registry{
		providers: make(map[string]Service),
		state:     state,
		alias:     alias,
		logger:    log,
	}
}

// Register associates a service with its provider name.
func (r *Registry) Register(provider string, svc Service) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Wrap the service with the validation decorator
	r.providers[provider] = NewValidationDecorator(svc, r.logger)
}

// Get retrieves the service for a given provider name.
func (r *Registry) Get(provider string) (Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	svc, ok := r.providers[provider]
	if !ok {
		return nil, fmt.Errorf("filesystem provider %s not found in registry", provider)
	}
	return svc, nil
}

func (r *Registry) RegisteredNames() ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names, nil
}

// DefaultProviderPrefix is the default provider to use when no prefix is specified.
const DefaultProviderPrefix = "onedrive"

// Cleanup removes all providers from the registry.
func (r *Registry) Cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers = make(map[string]Service)
}
