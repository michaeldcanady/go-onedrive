package infra

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

// MetadataCache defines the interface for caching drive item metadata.
// It supports retrieving, storing, and invalidating metadata to improve
// performance of file system operations.
type MetadataCache interface {
	// Get retrieves cached metadata for the given path.
	// It returns the metadata and a boolean indicating if the item was found.
	Get(ctx context.Context, path string) (*domain.Metadata, bool)
	// Put stores drive item metadata in the cache under the given path.
	Put(ctx context.Context, path string, m *domain.Metadata) error
	// Invalidate removes cached metadata for the given path.
	Invalidate(ctx context.Context, path string) error
}
