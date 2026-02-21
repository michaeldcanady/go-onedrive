package file

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
)

// MetadataCache defines the interface for caching drive item metadata.
// It supports retrieving, storing, and invalidating metadata to improve
// performance of file system operations.
type MetadataCache interface {
	// Get retrieves cached metadata for the given ID.
	// It returns the metadata and a boolean indicating if the item was found.
	Get(ctx context.Context, id string) (*file.Metadata, bool)
	// Put stores drive item metadata in the cache.
	Put(ctx context.Context, m *file.Metadata) error
	// Invalidate removes cached metadata for the given ID.
	Invalidate(ctx context.Context, id string) error
}
