package profile

import (
	"context"
	"fmt"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal2/core/shared"
)

// MemoryService is an in-memory implementation of the Profile Service.
type MemoryService struct {
	mu       sync.RWMutex
	profiles map[string]shared.Profile
}

// NewMemoryService initializes a new instance of MemoryService.
func NewMemoryService() *MemoryService {
	return &MemoryService{
		profiles: make(map[string]shared.Profile),
	}
}

// Get returns the profile with the specified name.
func (s *MemoryService) Get(ctx context.Context, name string) (shared.Profile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.profiles[name]
	if !ok {
		return shared.Profile{}, fmt.Errorf("profile %s not found", name)
	}
	return p, nil
}

// List returns a slice of all registered profiles.
func (s *MemoryService) List(ctx context.Context) ([]shared.Profile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]shared.Profile, 0, len(s.profiles))
	for _, p := range s.profiles {
		list = append(list, p)
	}
	return list, nil
}

// Create initializes and stores a new profile.
func (s *MemoryService) Create(ctx context.Context, name string) (shared.Profile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.profiles[name]; exists {
		return shared.Profile{}, fmt.Errorf("profile %s already exists", name)
	}
	p := shared.Profile{Name: name}
	s.profiles[name] = p
	return p, nil
}

// Delete removes the specified profile from storage.
func (s *MemoryService) Delete(ctx context.Context, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.profiles, name)
	return nil
}

// Exists checks if the given profile name is already in use.
func (s *MemoryService) Exists(ctx context.Context, name string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.profiles[name]
	return ok, nil
}
