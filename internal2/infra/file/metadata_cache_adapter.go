package file

import (
	"context"
	"encoding/json"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
)

type MetadataCacheAdapter struct {
	cache *abstractions.Cache2
}

func NewMetadataCacheAdapter(cache *abstractions.Cache2) *MetadataCacheAdapter {
	return &MetadataCacheAdapter{
		cache: cache,
	}
}

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

func (c *MetadataCacheAdapter) Put(ctx context.Context, m *file.Metadata) error {
	return c.cache.Set(
		ctx,
		func() ([]byte, error) { return []byte(m.Path), nil },
		func() ([]byte, error) { return json.Marshal(m) },
	)
}

func (c *MetadataCacheAdapter) Invalidate(ctx context.Context, path string) error {
	return c.cache.Delete(
		ctx,
		func() ([]byte, error) { return []byte(path), nil },
	)
}
