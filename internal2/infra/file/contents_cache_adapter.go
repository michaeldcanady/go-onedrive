package file

import (
	"context"

	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
)

var _ ContentsCache = (*ContentsCacheAdapter)(nil)

// ContentsCacheAdapter implements the [ContentsCache] interface using a generic
// cache implementation. It handles the serialization and deserialization of
// file contents using JSON.
type ContentsCacheAdapter struct {
	cache domaincache.Cache[file.Contents]
}

// NewContentsCacheAdapter constructs a new [ContentsCacheAdapter] using the
// provided cache implementation.
func NewContentsCacheAdapter(cache domaincache.Cache[file.Contents]) *ContentsCacheAdapter {
	return &ContentsCacheAdapter{
		cache: cache,
	}
}

// Get retrieves a *[file.Contents] value from the cache using the provided path
// as the lookup key.
func (c *ContentsCacheAdapter) Get(ctx context.Context, path string) (*file.Contents, bool) {
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

// Put stores a *[file.Contents] value in the cache under the given path.
func (c *ContentsCacheAdapter) Put(ctx context.Context, path string, m *file.Contents) error {
	if m == nil {
		return nil
	}
	return c.cache.Set(ctx, path, *m)
}
