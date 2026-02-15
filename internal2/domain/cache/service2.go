package cache

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
)

type Service2 interface {
	CreateCache(context.Context, string, func() abstractions.KeyValueStore) *abstractions.Cache2
	GetCache(context.Context, string) (*abstractions.Cache2, bool)
}
