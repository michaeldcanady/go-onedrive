package memory

import (
	"context"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/core"
)

type Store struct {
	mu    sync.RWMutex
	store map[string][]byte
}

func NewStore() *Store {
	return &Store{
		store: make(map[string][]byte),
	}
}

func (m *Store) Get(ctx context.Context, key []byte) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	m.mu.RLock()
	v, ok := m.store[string(key)]
	m.mu.RUnlock()

	if !ok {
		return nil, core.ErrKeyNotFound
	}

	// Copy to avoid exposing internal slice
	out := append([]byte(nil), v...)
	return out, nil
}

func (m *Store) Set(ctx context.Context, key, value []byte) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	// Copy to avoid aliasing user-provided slice
	k := string(append([]byte(nil), key...))
	v := append([]byte(nil), value...)

	m.mu.Lock()
	m.store[k] = v
	m.mu.Unlock()

	return nil
}

func (m *Store) Delete(ctx context.Context, key []byte) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	m.mu.Lock()
	delete(m.store, string(key))
	m.mu.Unlock()

	return nil
}
