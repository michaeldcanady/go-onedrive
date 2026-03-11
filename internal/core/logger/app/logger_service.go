package app

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"path/filepath"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
)

var _ domain.LoggerService = (*LoggerService)(nil)

type LoggerService struct {
	loggers  map[string]domain.Logger
	registry map[domain.Type]domain.LoggerProvider
	config   *domain.Options
}

// RegisterProvider implements [domain.LoggerService].
func (s *LoggerService) RegisterProvider(ype domain.Type, factory domain.LoggerProvider) {
	if s.registry == nil {
		s.registry = make(map[domain.Type]domain.LoggerProvider)
	}
	s.registry[ype] = factory
}

func NewLoggerService(opts ...domain.Option) (*LoggerService, error) {

	config := domain.NewOptions()

	if err := config.Apply(opts...); err != nil {
		return nil, errors.Join(errors.New("unable to apply logger options"), err)
	}

	return &LoggerService{
		loggers:  map[string]domain.Logger{},
		registry: map[domain.Type]domain.LoggerProvider{},
		config:   config,
	}, nil
}

func (s *LoggerService) buildLoggerOptions(id string) (domain.LoggerOptions, error) {
	config := domain.NewLoggerOptions()

	opts := []domain.LoggerOption{domain.WithLoggerLogLevel(s.config.LogLevel)}

	switch s.config.OutputDestination {
	case domain.OutputDestinationFile:
		opts = append(opts, domain.WithLoggerOutputDestinationFile(s.toPath(id)))
	case domain.OutputDestinationStandardOut:
		opts = append(opts, domain.WithLoggerOutputDestinationStandardOut())
	case domain.OutputDestinationStandardError:
		opts = append(opts, domain.WithLoggerOutputDestinationStandardError())
	}

	if err := config.Apply(opts...); err != nil {
		return domain.LoggerOptions{}, errors.Join(errors.New("unable to apply logger options"), err)
	}

	return *config, nil
}

// CreateLogger creates a new logger using the current factory.
func (s *LoggerService) CreateLogger(id string) (domain.Logger, error) {

	provider, ok := s.registry[s.config.Type]
	if !ok {
		return nil, fmt.Errorf("no logger provider registered for type: %v", s.config.Type)
	}

	config, err := s.buildLoggerOptions(id)
	if err != nil {
		return nil, errors.Join(errors.New("unable to build logger options"), err)
	}

	l, err := provider.Logger(config)
	if err != nil {
		return nil, err
	}

	s.loggers[id] = l

	return l, nil
}

func (s *LoggerService) GetContextLogger(ctx context.Context, name string) (domain.Logger, error) {
	l, err := s.GetLogger(name)
	if err != nil {
		return nil, err
	}
	return l.WithContext(ctx), nil
}

func (s *LoggerService) SetAllLevel(level string) {
	s.config.LogLevel = level
	maps.Keys(s.loggers)(func(key string) bool {
		s.SetLevel(key, level)
		return true
	})
}

func (s *LoggerService) GetLogger(id string) (domain.Logger, error) {
	l, ok := s.loggers[id]
	if !ok {
		return nil, ErrUnknownLogger
	}
	return l, nil
}

func (s *LoggerService) SetLevel(id, level string) error {
	l, err := s.GetLogger(id)
	if err != nil {
		return errors.Join(errors.New("unable to set level"), err)
	}

	l.SetLevel(level)
	return nil
}

// toPath combines the provided id and logsHome to create the log's path.
func (s *LoggerService) toPath(id string) string {
	return filepath.Join(s.config.LogsHome, fmt.Sprintf("%s.log", id))
}
