package storage

import (
	"fmt"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Service manages the lifecycle of storage backends.
type Service interface {
	Open(path string) (*bolt.DB, error)
	Close() error
}

// DefaultService implements the Service interface for bbolt.
type DefaultService struct {
	dbs map[string]*bolt.DB
	mu  sync.RWMutex
}

// NewDefaultService initializes a new storage service.
func NewDefaultService() *DefaultService {
	return &DefaultService{
		dbs: make(map[string]*bolt.DB),
	}
}

// Open opens the bbolt database at the specified path or returns an existing one.
func (s *DefaultService) Open(path string) (*bolt.DB, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if db, ok := s.dbs[path]; ok {
		return db, nil
	}

	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open database at %s: %w", path, err)
	}
	s.dbs[path] = db
	return db, nil
}

// Close closes all managed database connections.
func (s *DefaultService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var errs []error
	for path, db := range s.dbs {
		if err := db.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close database at %s: %w", path, err))
		}
		delete(s.dbs, path)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing databases: %v", errs)
	}
	return nil
}
