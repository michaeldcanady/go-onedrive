package infra

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	pkgcache "github.com/michaeldcanady/go-onedrive/pkg/cache"
)

var _ ContentsCache = (*ContentsCacheAdapter)(nil)

// ContentsCacheAdapter implements the [ContentsCache] interface using a generic
// cache implementation. It handles the serialization and deserialization of
// file contents using JSON.
type ContentsCacheAdapter struct {
	cache pkgcache.Cache[domain.Contents]
}

// NewContentsCacheAdapter constructs a new [ContentsCacheAdapter] using the
// provided cache implementation.
func NewContentsCacheAdapter(cache pkgcache.Cache[domain.Contents]) *ContentsCacheAdapter {
	return &ContentsCacheAdapter{
		cache: cache,
	}
}

// Get retrieves a *[domain.Contents] value from the cache using the provided path
// as the lookup key.
func (c *ContentsCacheAdapter) Get(ctx context.Context, path string) (*domain.Contents, bool) {
	m, err := c.cache.Get(ctx, path)
	if err != nil {
		return nil, false
	}
	return &m, true
}

// Invalidate removes the cached entry associated with the given path.
func (c *ContentsCacheAdapter) Invalidate(ctx context.Context, path string) error {
	return c.cache.Delete(ctx, path)
}

// Put stores a *[domain.Contents] value in the cache under the given path.
func (c *ContentsCacheAdapter) Put(ctx context.Context, path string, m *domain.Contents) error {
	if m == nil {
		return nil
	}
	return c.cache.Set(ctx, path, *m)
}
