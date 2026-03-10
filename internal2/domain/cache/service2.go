package cache

import (
	"context"
)

type Service2 interface {
	CreateCache(ctx context.Context, name string, storeFactory func() KeyValueStore) *Store
	GetCache(ctx context.Context, name string) (*Store, bool)
}
