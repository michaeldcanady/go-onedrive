package file

import (
	"context"
	"encoding/json"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
)

type cache interface {
	Delete(ctx context.Context, keySerializer abstractions.Serializer2) error
	Get(ctx context.Context, keySerializer abstractions.Serializer2, valueDeserializer abstractions.Deserializer2) error
	Set(ctx context.Context, keySerializer abstractions.Serializer2, valueSerializer abstractions.Serializer2) error
}

var _ ContentsCache = (*ContentsCacheAdapter)(nil)

type ContentsCacheAdapter struct {
	cache cache
}

func NewContentsCacheAdapter(cache cache) *ContentsCacheAdapter {
	return &ContentsCacheAdapter{
		cache: cache,
	}
}

// Get implements [ContentsCache].
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

// Invalidate implements [ContentsCache].
func (c *ContentsCacheAdapter) Invalidate(ctx context.Context, path string) error {
	return c.cache.Delete(
		ctx,
		func() ([]byte, error) { return []byte(path), nil },
	)
}

// Put implements [ContentsCache].
func (c *ContentsCacheAdapter) Put(ctx context.Context, path string, m *file.Contents) error {
	return c.cache.Set(
		ctx,
		func() ([]byte, error) { return []byte(path), nil },
		func() ([]byte, error) { return json.Marshal(m) },
	)
}
