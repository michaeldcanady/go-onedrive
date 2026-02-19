package file

import (
	"context"
	"encoding/json"
)

type MetadataListCacheAdapter struct {
	cache cache
}

func NewMetadataListCacheAdapter(cache cache) *MetadataListCacheAdapter {
	return &MetadataListCacheAdapter{
		cache: cache,
	}
}

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

func (c *MetadataListCacheAdapter) Put(ctx context.Context, path string, listing *Listing) error {
	return c.cache.Set(
		ctx,
		func() ([]byte, error) { return []byte(path), nil },
		func() ([]byte, error) { return json.Marshal(listing) },
	)
}

func (c *MetadataListCacheAdapter) Invalidate(ctx context.Context, path string) error {
	return c.cache.Delete(
		ctx,
		func() ([]byte, error) { return []byte(path), nil },
	)
}
