package cache

import (
	"context"
	"encoding/json"
	"errors"

	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/bolt"
)

var _ domaincache.Cache[any] = (*TypedCache[any])(nil)

type TypedCache[T any] struct {
	cache      *domaincache.Store
	serializer domaincache.SerializerDeserializer[T]
}

func NewTypedCache[T any](cache *domaincache.Store, serializer domaincache.SerializerDeserializer[T]) *TypedCache[T] {
	return &TypedCache[T]{
		cache:      cache,
		serializer: serializer,
	}
}

func (c *TypedCache[T]) Get(ctx context.Context, key string) (T, error) {
	var v T
	err := c.cache.Get(ctx,
		func() ([]byte, error) { return json.Marshal(key) },
		func(data []byte) error {
			var err error
			v, err = c.serializer.Deserialize(data)
			return err
		},
	)
	if errors.Is(err, bolt.ErrKeyNotFound) {
		return v, err
	}
	return v, err
}

func (c *TypedCache[T]) Set(ctx context.Context, key string, value T) error {
	return c.cache.Set(ctx,
		func() ([]byte, error) { return json.Marshal(key) },
		func() ([]byte, error) { return c.serializer.Serialize(value) },
	)
}

func (c *TypedCache[T]) Delete(ctx context.Context, key string) error {
	return c.cache.Delete(ctx, func() ([]byte, error) { return json.Marshal(key) })
}

func (c *TypedCache[T]) List(ctx context.Context, callback func(key string, value T) error) error {
	return c.cache.List(ctx, func(keyBytes []byte, valueBytes []byte) error {
		var key string
		if err := json.Unmarshal(keyBytes, &key); err != nil {
			return err
		}
		value, err := c.serializer.Deserialize(valueBytes)
		if err != nil {
			return err
		}
		return callback(key, value)
	})
}
