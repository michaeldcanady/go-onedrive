package infra

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	pkgcache "github.com/michaeldcanady/go-onedrive/pkg/cache"
)

// MetadataCacheAdapter implements the MetadataCache interface using a generic
// cache implementation. It handles the serialization and deserialization of
// drive item metadata using JSON.
type MetadataCacheAdapter struct {
	cache pkgcache.Cache[domain.Metadata]
}

// NewMetadataCacheAdapter constructs a new [MetadataCacheAdapter] using the
// provided cache implementation.
func NewMetadataCacheAdapter(cache pkgcache.Cache[domain.Metadata]) *MetadataCacheAdapter {
	return &MetadataCacheAdapter{
		cache: cache,
	}
}

// Get retrieves a *[domain.Metadata] value from the cache using the provided path
// as the lookup key.
func (c *MetadataCacheAdapter) Get(ctx context.Context, path string) (*domain.Metadata, bool) {
	m, err := c.cache.Get(ctx, path)
	if err != nil {
		return nil, false
	}
	return &m, true
}

// Put stores a *[domain.Metadata] value in the cache. The item's path is used
// as the cache key.
func (c *MetadataCacheAdapter) Put(ctx context.Context, path string, m *domain.Metadata) error {
	if m == nil {
		return nil
	}
	return c.cache.Set(ctx, path, *m)
}

// Invalidate removes the cached entry associated with the given path.
func (c *MetadataCacheAdapter) Invalidate(ctx context.Context, path string) error {
	return c.cache.Delete(ctx, path)
}
