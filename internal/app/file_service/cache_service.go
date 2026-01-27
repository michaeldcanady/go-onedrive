package driveservice

import (
	"context"
)

type CacheService interface {
	GetDrive(ctx context.Context, name string) (CachedChildren, error)
	SetDrive(ctx context.Context, name string, record CachedChildren) error
	GetItem(ctx context.Context, name string) (CachedItem, error)
	SetItem(ctx context.Context, name string, record CachedItem) error
}
