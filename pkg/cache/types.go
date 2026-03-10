package cache

import (
	"context"
	"errors"
)

var (
	// ErrKeyNotFound is returned when a requested key does not exist in the store.
	ErrKeyNotFound = errors.New("key not found")
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
