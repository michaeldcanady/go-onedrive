package file

import (
	"context"
	"encoding/json"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
)

// MetadataCacheAdapter implements the MetadataCache interface using a generic
// cache implementation. It handles the serialization and deserialization of
// drive item metadata using JSON.
type MetadataCacheAdapter struct {
	cache cache
}

// NewMetadataCacheAdapter constructs a new [MetadataCacheAdapter] using the
// provided cache implementation.
func NewMetadataCacheAdapter(cache cache) *MetadataCacheAdapter {
	return &MetadataCacheAdapter{
		cache: cache,
	}
}

// Get retrieves a *[file.Metadata] value from the cache using the provided path
// as the lookup key.
func (c *MetadataCacheAdapter) Get(ctx context.Context, path string) (*file.Metadata, bool) {
	var m file.Metadata

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

// Put stores a *[file.Metadata] value in the cache. The item's path is used
// as the cache key.
func (c *MetadataCacheAdapter) Put(ctx context.Context, m *file.Metadata) error {
	return c.cache.Set(
		ctx,
		func() ([]byte, error) { return []byte(m.Path), nil },
		func() ([]byte, error) { return json.Marshal(m) },
	)
}

// Invalidate removes the cached entry associated with the given path.
func (c *MetadataCacheAdapter) Invalidate(ctx context.Context, path string) error {
	return c.cache.Delete(
		ctx,
		func() ([]byte, error) { return []byte(path), nil },
	)
}
