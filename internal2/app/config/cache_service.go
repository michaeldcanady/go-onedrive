package config

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
)

// CacheService defines the interface used by the configuration Service
// to store and retrieve cached configuration data.
//
// Implementations must return core.ErrKeyNotFound when a configuration
// is not present in the cache.
type CacheService interface {
	GetConfiguration(ctx context.Context, name string) (config.Configuration3, error)
	SetConfiguration(ctx context.Context, name string, record config.Configuration3) error
}
