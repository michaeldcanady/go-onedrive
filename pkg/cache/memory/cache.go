package memory

import (
	"context"
	"sync"

	"github.com/michaeldcanady/go-onedrive/pkg/cache"
)

type Cache[K comparable, V any] struct {
	mu    sync.RWMutex
	store map[K]*cache.Entry[K, V]
}

func New[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		mu:    sync.RWMutex{},
		store: map[K]*cache.Entry[K, V]{},
	}
}

func (c *Cache[K, V]) GetEntry(ctx context.Context, key K) (*cache.Entry[K, V], error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.store[key]
	if !ok {
		return nil, cache.ErrKeyNotFound
	}

	return entry, nil
}

func (c *Cache[K, V]) NewEntry(ctx context.Context, key K) (*cache.Entry[K, V], error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var zero V
	return cache.NewEntry(key, zero), nil
}

func (c *Cache[K, V]) Remove(key K) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)

	return nil
}

func (c *Cache[K, V]) SetEntry(ctx context.Context, entry *cache.Entry[K, V]) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[entry.GetKey()] = entry

	return nil
}

func (c *Cache[K, V]) KeySerializer() cache.Serializer[K] {
	return nil
}

func (c *Cache[K, V]) Clear(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	c.store = map[K]*cache.Entry[K, V]{}
	return nil
}
