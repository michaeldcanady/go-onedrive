package cache

import (
	"context"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type Service2 struct {
	mu       sync.RWMutex
	logger   logging.Logger
	registry map[string]*abstractions.Cache2
}

func NewService2(logger logging.Logger) *Service2 {
	return &Service2{
		registry: make(map[string]*abstractions.Cache2),
		logger:   logger,
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

func (s *Service2) CreateCache(ctx context.Context, name string, storeFactory func() abstractions.KeyValueStore) *abstractions.Cache2 {
	correlationID := util.CorrelationIDFromContext(ctx)

	logger := s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("cache_name", name),
	)

	logger.Info(
		"creating cache instance",
		logging.String("event", eventCacheCreateStart),
	)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.registry[name]; exists {
		logger.Warn("cache already exists",
			logging.String("event", eventCacheCreateExists),
		)
		return s.registry[name]
	}

	cache := abstractions.NewCache2(storeFactory())
	s.registry[name] = cache

	logger.Info("cache created successfully",
		logging.String("event", eventCacheCreateFinish),
	)

	return cache
}

func (s *Service2) GetCache(ctx context.Context, name string) (*abstractions.Cache2, bool) {
	correlationID := util.CorrelationIDFromContext(ctx)

	logger := s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("cache_name", name),
	)

	logger.Debug("retrieving cache instance",
		logging.String("event", eventCacheGetStart),
	)

	s.mu.RLock()
	defer s.mu.RUnlock()

	cache, exists := s.registry[name]
	if exists {
		logger.Debug("cache found",
			logging.String("event", eventCacheGetHit),
		)
		return cache, true
	}

	logger.Debug("cache not found",
		logging.String("event", eventCacheGetMiss),
	)

	return nil, false
}
