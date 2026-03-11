package app

import (
	"context"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
	logger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	pkgcache "github.com/michaeldcanady/go-onedrive/pkg/cache"
)

type Service2 struct {
	mu       sync.RWMutex
	log      logger.Logger
	registry map[string]*pkgcache.Store
}

func NewService2(log logger.Logger) *Service2 {
	return &Service2{
		registry: make(map[string]*pkgcache.Store),
		log:      log,
	}
}

// ───────────────────────────────────────────────────────────────────────────────
// Event Taxonomy (cachedomain.registry)
// ───────────────────────────────────────────────────────────────────────────────

const (
	eventCacheCreateStart  = "cachedomain.registry.create.start"
	eventCacheCreateFinish = "cachedomain.registry.create.finish"
	eventCacheCreateExists = "cachedomain.registry.create.exists"

	eventCacheGetStart = "cachedomain.registry.get.start"
	eventCacheGetHit   = "cachedomain.registry.get.hit"
	eventCacheGetMiss  = "cachedomain.registry.get.miss"
)

// CreateCache instantiates a new cache store and adds it to the registry.
func (s *Service2) CreateCache(ctx context.Context, name string, storeFactory func() pkgcache.KeyValueStore) *pkgcache.Store {
	correlationID := util.CorrelationIDFromContext(ctx)

	log := s.log.WithContext(ctx).With(
		logger.String("correlation_id", correlationID),
		logger.String("cache_name", name),
	)

	log.Info(
		"creating cache instance",
		logger.String("event", eventCacheCreateStart),
	)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.registry[name]; exists {
		log.Warn("cache already exists",
			logger.String("event", eventCacheCreateExists),
		)
		return s.registry[name]
	}

	cache := pkgcache.NewStore(storeFactory())
	s.registry[name] = cache

	log.Info("cache created successfully",
		logger.String("event", eventCacheCreateFinish),
	)

	return cache
}

// GetCache returns the cache with the provided name.
func (s *Service2) GetCache(ctx context.Context, name string) (*pkgcache.Store, bool) {
	correlationID := util.CorrelationIDFromContext(ctx)

	log := s.log.WithContext(ctx).With(
		logger.String("correlation_id", correlationID),
		logger.String("cache_name", name),
	)

	log.Debug("retrieving cache instance",
		logger.String("event", eventCacheGetStart),
	)

	s.mu.RLock()
	defer s.mu.RUnlock()

	cache, exists := s.registry[name]
	if exists {
		log.Debug("cache found",
			logger.String("event", eventCacheGetHit),
		)
		return cache, true
	}

	log.Debug("cache not found",
		logger.String("event", eventCacheGetMiss),
	)

	return nil, false
}
