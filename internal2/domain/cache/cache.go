package cache

import (
	"context"
)

// Cache is a generic interface for a key-value store.
type Cache[T any] interface {
	Get(ctx context.Context, key string) (T, error)
	Set(ctx context.Context, key string, value T) error
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, callback func(key string, value T) error) error
}
