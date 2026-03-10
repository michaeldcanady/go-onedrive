package cache

import (
	"context"
	"errors"
)

// Store is a middle-level byte-based cache implementation.
type Store struct {
	store KeyValueStore
}

func NewStore(store KeyValueStore) *Store {
	return &Store{store: store}
}

func (c *Store) Get(
	ctx context.Context,
	keySerializer SerializerFunc,
	valueDeserializer DeserializerFunc,
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

func (c *Store) Set(
	ctx context.Context,
	keySerializer SerializerFunc,
	valueSerializer SerializerFunc,
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

func (c *Store) Delete(
	ctx context.Context,
	keySerializer SerializerFunc,
) error {
	keyBytes, err := keySerializer()
	if err != nil {
		return err
	}
	return c.store.Delete(ctx, keyBytes)
}

func (c *Store) List(
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
