package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.etcd.io/bbolt"
)

// BoltDB is a type alias for the underlying bbolt database instance.
type BoltDB = *bbolt.DB

// Service provides access to the underlying persistent database.
// It ensures that the database is correctly initialized and safe for concurrent use.
type Service[T any] interface {
	// DB returns the underlying database instance.
	DB() T

	// Close releases all resources associated with the storage service.
	// Subsequent calls to DB() after Close() will return nil.
	Close() error
}

type storageService[T any] struct {
	db        T
	mu        sync.Mutex
	path      string
	closeFunc func(T) error
}

// NewStorageService returns a new [Service] that manages a [bbolt] database at the specified path.
// It automatically creates the parent directory if it does not exist.
func NewStorageService(dbPath string) (Service[*bbolt.DB], error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open bbolt database: %w", err)
	}

	return &storageService[*bbolt.DB]{
		db:   db,
		path: dbPath,
		closeFunc: func(db *bbolt.DB) error {
			return db.Close()
		},
	}, nil
}

// NewService returns a new [Service] with the provided database and close function.
func NewService[T any](db T, closeFunc func(T) error) Service[T] {
	return &storageService[T]{
		db:        db,
		closeFunc: closeFunc,
	}
}

func (s *storageService[T]) DB() T {
	return s.db
}

func (s *storageService[T]) Shutdown(ctx context.Context) error {
	return s.Close()
}

func (s *storageService[T]) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closeFunc != nil {
		err := s.closeFunc(s.db)
		var zero T
		s.db = zero
		s.closeFunc = nil
		return err
	}
	return nil
}
