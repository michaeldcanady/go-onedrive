package registry

import (
	"fmt"
	"sync"

	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

type Registry struct {
	mu        sync.RWMutex
	providers map[string]domainfs.Service
}

func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]domainfs.Service),
	}
}

func (r *Registry) Register(name string, provider domainfs.Service) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[name] = provider
}

func (r *Registry) Get(name string) (domainfs.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return p, nil
}
