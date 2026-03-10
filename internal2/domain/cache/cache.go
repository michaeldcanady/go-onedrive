package cache

import (
	"context"
	"errors"
)

// Cache is a generic interface for a key-value store.
type Cache[T any] interface {
	Get(ctx context.Context, key string) (T, error)
	Set(ctx context.Context, key string, value T) error
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, callback func(key string, value T) error) error
}

// Entry represents a single cache entry.
type Entry[K, V any] struct {
	key   K
	value V
}

func NewEntry[K, V any](key K, value V) *Entry[K, V] {
	return &Entry[K, V]{
		key:   key,
		value: value,
	}
}

func (c *Entry[K, V]) GetKey() K {
	return c.key
}

func (c *Entry[K, V]) GetValue() V {
	return c.value
}

func (c *Entry[K, V]) SetValue(value V) {
	c.value = value
}

// SerializerFunc is a function that returns a byte representation of a value.
type SerializerFunc func() ([]byte, error)

// DeserializerFunc is a function that populates a value from a byte representation.
type DeserializerFunc func([]byte) error

// Serializer is a generic interface for serializing types.
type Serializer[T any] interface {
	Serialize(T) ([]byte, error)
}

// Deserializer is a generic interface for deserializing types.
type Deserializer[T any] interface {
	Deserialize([]byte) (T, error)
}

// SerializerDeserializer is a generic interface for both serializing and deserializing types.
type SerializerDeserializer[T any] interface {
	Serializer[T]
	Deserializer[T]
}

// KeyValueStore is a low-level byte-based key-value store interface.
type KeyValueStore interface {
	Get(context.Context, []byte) ([]byte, error)
	Set(context.Context, []byte, []byte) error
	Delete(context.Context, []byte) error
	List(context.Context) ([][]byte, [][]byte, error)
}

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
