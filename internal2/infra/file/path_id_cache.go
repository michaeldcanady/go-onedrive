package file

import (
	"context"

	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
)

type PathIDCache interface {
	Get(ctx context.Context, path string) (string, bool)
	Put(ctx context.Context, path string, id string) error
	Invalidate(ctx context.Context, path string) error
}

type PathIDCacheAdapter struct {
	cache domaincache.Cache[string]
}

func NewPathIDCacheAdapter(cache domaincache.Cache[string]) *PathIDCacheAdapter {
	return &PathIDCacheAdapter{
		cache: cache,
	}
}

func (c *PathIDCacheAdapter) Get(ctx context.Context, path string) (string, bool) {
	id, err := c.cache.Get(ctx, path)
	if err != nil {
		return "", false
	}
	return id, true
}

func (c *PathIDCacheAdapter) Put(ctx context.Context, path string, id string) error {
	return c.cache.Set(ctx, path, id)
}

func (c *PathIDCacheAdapter) Invalidate(ctx context.Context, path string) error {
	return c.cache.Delete(ctx, path)
}
