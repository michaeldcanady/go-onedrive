package cache

import (
	"context"
	"sync"

	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type Service2 struct {
	mu       sync.RWMutex
	log      logger.Logger
	registry map[string]*domaincache.Store
}

func NewService2(log logger.Logger) *Service2 {
	return &Service2{
		registry: make(map[string]*domaincache.Store),
		log:      log,
	}
}

// ───────────────────────────────────────────────────────────────────────────────
// Event Taxonomy (cache.registry)
// ───────────────────────────────────────────────────────────────────────────────

const (
	eventCacheCreateStart  = "cache.registry.create.start"
	eventCacheCreateFinish = "cache.registry.create.finish"
	eventCacheCreateExists = "cache.registry.create.exists"

	eventCacheGetStart = "cache.registry.get.start"
	eventCacheGetHit   = "cache.registry.get.hit"
	eventCacheGetMiss  = "cache.registry.get.miss"
)

func (s *Service2) CreateCache(ctx context.Context, name string, storeFactory func() domaincache.KeyValueStore) *domaincache.Store {
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

	cache := domaincache.NewStore(storeFactory())
	s.registry[name] = cache

	log.Info("cache created successfully",
		logger.String("event", eventCacheCreateFinish),
	)

	return cache
}

func (s *Service2) GetCache(ctx context.Context, name string) (*domaincache.Store, bool) {
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
