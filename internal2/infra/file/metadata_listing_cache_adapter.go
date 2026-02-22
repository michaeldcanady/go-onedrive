package file

import (
	"context"
	"encoding/json"
)

// MetadataListCacheAdapter implements the ListingCache interface using a
// generic cache implementation. It handles the serialization and
// deserialization of folder listings (Metadata Listing) using JSON.
type MetadataListCacheAdapter struct {
	cache cache
}

// NewMetadataListCacheAdapter constructs a new [MetadataListCacheAdapter] using
// the provided cache implementation.
func NewMetadataListCacheAdapter(cache cache) *MetadataListCacheAdapter {
	return &MetadataListCacheAdapter{
		cache: cache,
	}
}

// Get retrieves a cached *[Listing] for the given path.
func (c *MetadataListCacheAdapter) Get(ctx context.Context, path string) (*Listing, bool) {
	var listing Listing

	err := c.cache.Get(
		ctx,
		func() ([]byte, error) { return []byte(path), nil },
		func(b []byte) error { return json.Unmarshal(b, &listing) },
	)

	if err != nil {
		return nil, false
	}
	return &listing, true
}

// Put stores a *[Listing] in the cache under the provided path.
func (c *MetadataListCacheAdapter) Put(ctx context.Context, path string, listing *Listing) error {
	return c.cache.Set(
		ctx,
		func() ([]byte, error) { return []byte(path), nil },
		func() ([]byte, error) { return json.Marshal(listing) },
	)
}

// Invalidate removes the cached listing associated with the given path.
func (c *MetadataListCacheAdapter) Invalidate(ctx context.Context, path string) error {
	return c.cache.Delete(
		ctx,
		func() ([]byte, error) { return []byte(path), nil },
	)
}
