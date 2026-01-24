package configurationservice

import (
	"context"
	"fmt"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal/cachev2/abstractions"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/event"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	"github.com/spf13/viper"
)

type Service struct {
	config       config.Configuration2
	cacheService abstractions.Cache[string, config.Configuration2]
	path         string
	lock         sync.RWMutex
	publisher    event.Publisher
	logger       logging.Logger
}

// New creates a new ConfigurationService instance.
func New(path string, publisher event.Publisher, cache abstractions.Cache[string, config.Configuration2], logger logging.Logger) *Service {
	return &Service{
		config:    nil,
		path:      path,
		lock:      sync.RWMutex{},
		publisher: publisher,
		logger:    logger,
	}
}

func (s *Service) LoadConfiguration(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	if s.path == "" {
		s.logger.Warn("configuration file path is empty; cannot load configuration")
		return fmt.Errorf("configuration file path is empty")
	}

	s.config = config.NewViperAdapter(viper.New())
	s.logger.Info("loading configuration from file", logging.String("path", s.path))
	s.config.SetConfigFile(s.path)
	if err := s.config.ReadInConfig(); err != nil {
		s.logger.Error("failed to read configuration", logging.Any("error", err))
		return fmt.Errorf("failed to read configuration: %w", err)
	}

	s.logger.Info("configuration loaded successfully")

	return nil
}

func (s *Service) WriteConfiguration(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.path == "" {
		s.logger.Warn("configuration file path is empty; cannot write configuration")
		return fmt.Errorf("configuration file path is empty")
	}

	s.logger.Info("writing configuration to file", logging.String("path", s.path))
	if err := s.config.WriteConfig(); err != nil {
		s.logger.Error("failed to write configuration", logging.Any("error", err))
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	s.logger.Info("configuration written successfully")

	return nil
}

func (s *Service) SetConfigFile(ctx context.Context, path string) {
	if ctx.Err() != nil {
		return
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	old := s.path

	s.path = path

	if s.publisher == nil {
		s.logger.Warn("event publisher is nil; skipping configuration file changed")
	} else {
		s.logger.Debug("publishing configuration file changed event", logging.String("path", path))
		if err := s.publisher.Publish(ctx, newConfigurationFileChangedEvent(old, path)); err != nil {
			s.logger.Warn("failed to publish configuration file changed event", logging.Any("error", err))
		}
	}
}

// GetString retrieves a string value from the configuration by key.
func (s *Service) GetString(ctx context.Context, key string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	if s.config == nil {
		s.logger.Warn("configuration source not defined")
		if err := s.LoadConfiguration(ctx); err != nil {
			return "", err
		}
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.config.GetString(key), nil
}

func (s *Service) GetStringDefault(ctx context.Context, key string, defaultValue string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	val, err := s.GetString(ctx, key)
	if err != nil {
		return "", err
	}
	if val == "" {
		val = defaultValue
	}
	return val, nil
}

// Get retrieves a value from the configuration by key.
func (s *Service) Get(ctx context.Context, key string) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if s.config == nil {
		s.logger.Warn("configuration source not defined")
		if err := s.LoadConfiguration(ctx); err != nil {
			return "", err
		}
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.config.Get(key), nil
}

// Set sets a value in the configuration by key.
func (s *Service) Set(ctx context.Context, key string, value any) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if s.config == nil {
		s.logger.Warn("configuration not loaded; cannot set value")
		if err := s.LoadConfiguration(ctx); err != nil {
			return err
		}
	}

	old, err := s.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("unable to retrieve current value: %w", err)
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.config.Set(key, value)

	if s.publisher == nil {
		s.logger.Warn("event publisher is nil; skipping configuration changed event publish")
	} else {
		s.logger.Debug("publishing configuration updated event", logging.String("key", key), logging.Any("value", value))
		if err := s.publisher.Publish(ctx, newConfigurationUpdatedEvent(key, old, value)); err != nil {
			s.logger.Warn("failed to publish configuration changed event", logging.Any("error", err))
		}
	}

	return nil
}

// SetString sets a string value in the configuration by key.
func (s *Service) SetString(ctx context.Context, key string, value string) error {

	return s.Set(ctx, key, value)
}
