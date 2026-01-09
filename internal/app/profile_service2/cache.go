package cacheservice

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/cache/fsstore"
)

var _ Cache = (*fsstore.FSStore)(nil)

type Cache interface {
	Get(context.Context, string) ([]byte, error)
	Put(context.Context, string, []byte) error
	Delete(context.Context, string) error
	Clear(context.Context) error
}
