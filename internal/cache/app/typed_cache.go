package app

import (
	"context"
	"encoding/json"

	domaincache "github.com/michaeldcanady/go-onedrive/internal/cache/domain"
)

var _ domaincache.Cache[any] = (*TypedCache[any])(nil)

type TypedCache[T any] struct {
	cache           *domaincache.Store
	keySerializer   domaincache.SerializerDeserializer[string]
	valueSerializer domaincache.SerializerDeserializer[T]
}

func NewTypedCache[T any](cache *domaincache.Store, valueSerializer domaincache.SerializerDeserializer[T]) *TypedCache[T] {
	return &TypedCache[T]{
		cache:           cache,
		keySerializer:   &JSONSerializerDeserializer[string]{},
		valueSerializer: valueSerializer,
	}
}

func (c *TypedCache[T]) Get(ctx context.Context, key string) (T, error) {
	var v T
	err := c.cache.Get(ctx,
		func() ([]byte, error) { return c.keySerializer.Serialize(key) },
		func(data []byte) error {
			var err error
			v, err = c.valueSerializer.Deserialize(data)
			return err
		},
	)
	return v, err
}

func (c *TypedCache[T]) Set(ctx context.Context, key string, value T) error {
	return c.cache.Set(ctx,
		func() ([]byte, error) { return c.keySerializer.Serialize(key) },
		func() ([]byte, error) { return c.valueSerializer.Serialize(value) },
	)
}

func (c *TypedCache[T]) Delete(ctx context.Context, key string) error {
	return c.cache.Delete(ctx, func() ([]byte, error) { return json.Marshal(key) })
}

func (c *TypedCache[T]) List(ctx context.Context, callback func(key string, value T) error) error {
	return c.cache.List(ctx, func(keyBytes []byte, valueBytes []byte) error {
		var key string
		key, err := c.keySerializer.Deserialize(keyBytes)
		if err != nil {
			return err
		}
		value, err := c.valueSerializer.Deserialize(valueBytes)
		if err != nil {
			return err
		}
		return callback(key, value)
	})
}
