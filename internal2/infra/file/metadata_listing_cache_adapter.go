package file

import (
	"context"
	"encoding/json"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
)

type metadataListCacheAdapter struct {
	cache *abstractions.Cache2
}

func (c *metadataListCacheAdapter) Get(ctx context.Context, path string) (*Listing, bool) {
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

func (c *metadataListCacheAdapter) Put(ctx context.Context, path string, listing *Listing) error {
	return c.cache.Set(
		ctx,
		func() ([]byte, error) { return []byte(path), nil },
		func() ([]byte, error) { return json.Marshal(listing) },
	)
}

func (c *metadataListCacheAdapter) Invalidate(ctx context.Context, path string) error {
	return c.cache.Delete(
		ctx,
		func() ([]byte, error) { return []byte(path), nil },
	)
}
