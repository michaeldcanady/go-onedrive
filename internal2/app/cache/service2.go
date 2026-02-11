package cache

import (
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
)

type Service2 struct {
	mu       sync.RWMutex
	registry map[string]*abstractions.Cache2
}

func NewService2() *Service2 {
	return &Service2{
		registry: make(map[string]*abstractions.Cache2),
	}
}

func (s *Service2) CreateCache(name string, storeFactory func() abstractions.KeyValueStore) *abstractions.Cache2 {
	s.mu.Lock()
	defer s.mu.Unlock()

	cache := abstractions.NewCache2(storeFactory())
	s.registry[name] = cache
	return cache
}

func (s *Service2) GetCache(name string) (*abstractions.Cache2, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cache, exists := s.registry[name]
	return cache, exists
}
