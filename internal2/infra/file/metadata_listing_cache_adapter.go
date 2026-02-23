package file

import (
	"context"

	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
)

// MetadataListCacheAdapter implements the ListingCache interface using a
// generic cache implementation. It handles the serialization and
// deserialization of folder listings (Metadata Listing) using JSON.
type MetadataListCacheAdapter struct {
	cache domaincache.Cache[Listing]
}

// NewMetadataListCacheAdapter constructs a new [MetadataListCacheAdapter] using
// the provided cache implementation.
func NewMetadataListCacheAdapter(cache domaincache.Cache[Listing]) *MetadataListCacheAdapter {
	return &MetadataListCacheAdapter{
		cache: cache,
	}
}

// Get retrieves a cached *[Listing] for the given path.
func (c *MetadataListCacheAdapter) Get(ctx context.Context, path string) (*Listing, bool) {
	listing, err := c.cache.Get(ctx, path)
	if err != nil {
		return nil, false
	}
	return &listing, true
}

// Put stores a *[Listing] in the cache under the provided path.
func (c *MetadataListCacheAdapter) Put(ctx context.Context, path string, listing *Listing) error {
	if listing == nil {
		return nil
	}
	return c.cache.Set(ctx, path, *listing)
}

// Invalidate removes the cached listing associated with the given path.
func (c *MetadataListCacheAdapter) Invalidate(ctx context.Context, path string) error {
	return c.cache.Delete(ctx, path)
}
