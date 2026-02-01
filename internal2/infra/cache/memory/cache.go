package memory

import (
	"context"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/core"
)

type Cache[K comparable, V any] struct {
	mu    sync.RWMutex
	store map[K]*abstractions.Entry[K, V]
}

func New[Entry *abstractions.Entry[K, V], K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		mu:    sync.RWMutex{},
		store: map[K]*abstractions.Entry[K, V]{},
	}
}

func (c *Cache[K, V]) GetEntry(ctx context.Context, key K) (*abstractions.Entry[K, V], error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.store[key]
	if !ok {
		return nil, core.ErrKeyNotFound
	}

	return entry, nil
}

func (c *Cache[K, V]) NewEntry(ctx context.Context, key K) (*abstractions.Entry[K, V], error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var zero V
	return abstractions.NewEntry(key, zero), nil
}

func (c *Cache[K, V]) Remove(key K) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)

	return nil
}

func (c *Cache[K, V]) SetEntry(ctx context.Context, entry *abstractions.Entry[K, V]) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[entry.GetKey()] = entry

	return nil
}

func (c *Cache[K, V]) KeySerializer() abstractions.Serializer[K] {
	return nil
}

func (c *Cache[K, V]) Clear(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	c.store = map[K]*abstractions.Entry[K, V]{}
	return nil
}
