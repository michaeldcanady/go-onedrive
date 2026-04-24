package storage

import (
	"fmt"
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
	db *bolt.DB
}

// NewDefaultService initializes a new storage service.
func NewDefaultService() *DefaultService {
	return &DefaultService{}
}

// Open opens the bbolt database at the specified path.
func (s *DefaultService) Open(path string) (*bolt.DB, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	s.db = db
	return db, nil
}

// Close closes the database connection.
func (s *DefaultService) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
