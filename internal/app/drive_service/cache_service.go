package driveservice

import (
	"context"
)

type CacheService interface {
	GetDrive(ctx context.Context, name string) (CachedChildren, error)
	SetDrive(ctx context.Context, name string, record CachedChildren) error
}
