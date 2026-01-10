package loggingservice

import (
	"context"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal/logging"
	"go.uber.org/zap"
)

type Service struct {
	logger  logging.Logger
	lock    sync.RWMutex
	level   string
	logHome string
}

// TODO: how to manage and change settings?

func New(level string, logHome string) *Service {
	return &Service{
		level:   level,
		logHome: logHome,
	}
}

func (s *Service) initializeLogger(ctx context.Context) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	cfg := zap.NewProductionConfig()

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	s.logger = logging.NewZapLoggerAdapter(logger)
	return nil
}

// Info logs a message at InfoLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (s *Service) Info(ctx context.Context, msg string, kv ...logging.Field) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if s.logger == nil {
		if err := s.initializeLogger(ctx); err != nil {
			return err
		}
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.logger.Info(msg, kv...)
	return nil
}

// Warn logs a message at WarnLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (s *Service) Warn(ctx context.Context, msg string, kv ...logging.Field) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if s.logger == nil {
		if err := s.initializeLogger(ctx); err != nil {
			return err
		}
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.logger.Warn(msg, kv...)
	return nil
}

// Error logs a message at ErrorLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (s *Service) Error(ctx context.Context, msg string, kv ...logging.Field) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if s.logger == nil {
		if err := s.initializeLogger(ctx); err != nil {
			return err
		}
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.logger.Error(msg, kv...)
	return nil
}

// Debug logs a message at DebugLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (s *Service) Debug(ctx context.Context, msg string, kv ...logging.Field) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if s.logger == nil {
		if err := s.initializeLogger(ctx); err != nil {
			return err
		}
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.logger.Debug(msg, kv...)
	return nil
}
