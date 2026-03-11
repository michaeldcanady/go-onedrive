package infra

import (
	"context"

	pkgcache "github.com/michaeldcanady/go-onedrive/pkg/cache"
)

type PathIDCache interface {
	Get(ctx context.Context, path string) (string, bool)
	Put(ctx context.Context, path string, id string) error
	Invalidate(ctx context.Context, path string) error
}

type PathIDCacheAdapter struct {
	cache pkgcache.Cache[string]
}

func NewPathIDCacheAdapter(cache pkgcache.Cache[string]) *PathIDCacheAdapter {
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
