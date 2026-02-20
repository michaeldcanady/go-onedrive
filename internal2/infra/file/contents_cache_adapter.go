package file

import (
	"context"
	"encoding/json"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
)

var _ ContentsCache = (*ContentsCacheAdapter)(nil)

type ContentsCacheAdapter struct {
	cache cache
}

// NewContentsCacheAdapter constructs a new [ContentsCacheAdapter] using the
// provided cache implementation.
func NewContentsCacheAdapter(cache cache) *ContentsCacheAdapter {
	return &ContentsCacheAdapter{
		cache: cache,
	}
}

// Get retrieves a *[file.Contents] value from the cache using the provided path
// as the lookup key.
func (c *ContentsCacheAdapter) Get(ctx context.Context, path string) (*file.Contents, bool) {
	var m file.Contents

	err := c.cache.Get(
		ctx,
		func() ([]byte, error) { return []byte(path), nil },
		func(b []byte) error { return json.Unmarshal(b, &m) },
	)

	if err != nil {
		return nil, false
	}
	return &m, true
}

// Invalidate removes the cached entry associated with the given path.
func (c *ContentsCacheAdapter) Invalidate(ctx context.Context, path string) error {
	return c.cache.Delete(
		ctx,
		func() ([]byte, error) { return []byte(path), nil },
	)
}

// Put stores a *[file.Contents] value in the cache under the given path.
func (c *ContentsCacheAdapter) Put(ctx context.Context, path string, m *file.Contents) error {
	return c.cache.Set(
		ctx,
		func() ([]byte, error) { return []byte(path), nil },
		func() ([]byte, error) { return json.Marshal(m) },
	)
}
