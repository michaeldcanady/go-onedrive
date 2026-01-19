package memory

import (
	"context"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal/cachev2/abstractions"
	"github.com/michaeldcanady/go-onedrive/internal/cachev2/core"
)

type Cache[Entry abstractions.Entry[K, V], K comparable, V any] struct {
	mu    sync.RWMutex
	store map[K]*Entry
}

func (c *Cache[Entry, K, V]) GetEntry(ctx context.Context, key K) (*Entry, error) {
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

func (c *Cache[K, V]) NewEntry(ctx context.Context, key K) (*Entry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var zero V
	return abstractions.NewEntry(key, zero), nil
}

func (c *Cache[K, V]) Remove(key K) error {
	c.mu.Lock()
	defer c.mu.Unlock()

}

func (c *Cache[K, V]) SetEntry(ctx context.Context, entry *Entry) error {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = entry

	return nil
}

func (c *Cache[K, V]) KeySerializer() abstractions.Serializer[K] {
	return nil
}
