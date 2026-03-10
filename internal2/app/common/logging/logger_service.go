package logging

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"path/filepath"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
)

var _ logger.LoggerService = (*LoggerService)(nil)

type LoggerService struct {
	loggers  map[string]logger.Logger
	registry map[logger.Type]logger.LoggerProvider
	config   *logger.Options
}

// RegisterProvider implements [logger.LoggerService].
func (s *LoggerService) RegisterProvider(ype logger.Type, factory logger.LoggerProvider) {
	if s.registry == nil {
		s.registry = make(map[logger.Type]logger.LoggerProvider)
	}
	s.registry[ype] = factory
}

func NewLoggerService(opts ...logger.Option) (*LoggerService, error) {

	config := logger.NewOptions()

	if err := config.Apply(opts...); err != nil {
		return nil, errors.Join(errors.New("unable to apply logger options"), err)
	}

	return &LoggerService{
		loggers:  map[string]logger.Logger{},
		registry: map[logger.Type]logger.LoggerProvider{},
		config:   config,
	}, nil
}

func (s *LoggerService) buildLoggerOptions(id string) (logger.LoggerOptions, error) {
	config := logger.NewLoggerOptions()

	opts := []logger.LoggerOption{logger.WithLoggerLogLevel(s.config.LogLevel)}

	switch s.config.OutputDestination {
	case logger.OutputDestinationFile:
		opts = append(opts, logger.WithLoggerOutputDestinationFile(s.toPath(id)))
	case logger.OutputDestinationStandardOut:
		opts = append(opts, logger.WithLoggerOutputDestinationStandardOut())
	case logger.OutputDestinationStandardError:
		opts = append(opts, logger.WithLoggerOutputDestinationStandardError())
	}

	if err := config.Apply(opts...); err != nil {
		return logger.LoggerOptions{}, errors.Join(errors.New("unable to apply logger options"), err)
	}

	return *config, nil
}

// createLogger creates a new logger using the current factory.
func (s *LoggerService) CreateLogger(id string) (logger.Logger, error) {

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

func (s *LoggerService) GetContextLogger(ctx context.Context, name string) (logger.Logger, error) {
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

func (s *LoggerService) GetLogger(id string) (logger.Logger, error) {
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
