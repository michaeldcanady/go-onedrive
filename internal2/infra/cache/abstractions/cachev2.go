package abstractions

import (
	"context"
	"errors"
)

type Cache2 struct {
	store KeyValueStore
}

func NewCache2(store KeyValueStore) *Cache2 {
	return &Cache2{store: store}
}

func (c *Cache2) Get(
	ctx context.Context,
	keySerializer Serializer2,
	valueDeserializer Deserializer2,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	keyBytes, err := keySerializer()
	if err != nil {
		return err
	}

	if len(keyBytes) == 0 {
		return errors.New("key is empty")
	}

	raw, err := c.store.Get(ctx, keyBytes)
	if err != nil {
		return err
	}

	return valueDeserializer(raw)
}

func (c *Cache2) Set(
	ctx context.Context,
	keySerializer Serializer2,
	valueSerializer Serializer2,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	keyBytes, err := keySerializer()
	if err != nil {
		return err
	}

	valueBytes, err := valueSerializer()
	if err != nil {
		return err
	}

	return c.store.Set(ctx, keyBytes, valueBytes)
}

func (c *Cache2) Delete(
	ctx context.Context,
	keySerializer Serializer2,
) error {
	keyBytes, err := keySerializer()
	if err != nil {
		return err
	}
	return c.store.Delete(ctx, keyBytes)
}

func (c *Cache2) List(
	ctx context.Context,
	callback func(key []byte, value []byte) error,
) error {
	keys, values, err := c.store.List(ctx)
	if err != nil {
		return err
	}

	for i := range keys {
		if err := callback(keys[i], values[i]); err != nil {
			return err
		}
	}

	return nil
}
