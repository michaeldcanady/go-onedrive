package infra

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/fs/shared/domain"
)

// ContentsCache defines the interface for caching file contents. It allows
// for retrieving, storing, and invalidating cached file data to improve
// performance and reduce network overhead.
type ContentsCache interface {
	// Get retrieves cached file contents for the given ID.
	// It returns the contents and a boolean indicating if the item was found.
	Get(ctx context.Context, id string) (*domain.Contents, bool)
	// Put stores file contents in the cache for the given ID.
	Put(ctx context.Context, id string, m *domain.Contents) error
	// Invalidate removes cached file contents for the given ID.
	Invalidate(ctx context.Context, id string) error
}
