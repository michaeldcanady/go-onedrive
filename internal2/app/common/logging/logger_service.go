package logging

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"path/filepath"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
)

var _ domainlogger.LoggerService = (*LoggerService)(nil)

type LoggerService struct {
	loggers  map[string]logging.Logger
	registry map[logging.Type]domainlogger.LoggerProvider
	config   *logger.Options
}

// RegisterProvider implements [logger.LoggerService].
func (s *LoggerService) RegisterProvider(ype logging.Type, factory domainlogger.LoggerProvider) {
	if s.registry == nil {
		s.registry = make(map[logging.Type]domainlogger.LoggerProvider)
	}
	s.registry[ype] = factory
}

func NewLoggerService(opts ...logger.Option) (*LoggerService, error) {

	config := logger.NewOptions()

	if err := config.Apply(opts...); err != nil {
		return nil, errors.Join(errors.New("unable to apply logger options"), err)
	}

	return &LoggerService{
		loggers:  map[string]logging.Logger{},
		registry: map[logging.Type]domainlogger.LoggerProvider{},
		config:   config,
	}, nil
}

func (s *LoggerService) buildLoggerOptions(id string) (logging.LoggerOptions, error) {
	config := logging.NewLoggerOptions()

	opts := []logging.LoggerOption{logging.WithLogLevel(s.config.LogLevel)}

	switch s.config.OutputDestination {
	case logging.OutputDestinationFile:
		opts = append(opts, logging.WithOutputDestinationFile(s.toPath(id)))
	case logging.OutputDestinationStandardOut:
		opts = append(opts, logging.WithOutputDestinationStandardOut())
	case logging.OutputDestinationStandardError:
		opts = append(opts, logging.WithOutputDestinationStandardError())
	}

	if err := config.Apply(opts...); err != nil {
		return logging.LoggerOptions{}, errors.Join(errors.New("unable to apply logger options"), err)
	}

	return *config, nil
}

// createLogger creates a new logger using the current factory.
func (s *LoggerService) CreateLogger(id string) (logging.Logger, error) {

	provider, ok := s.registry[s.config.Type]
	if !ok {
		return nil, fmt.Errorf("no logger provider registered for type: %v", s.config.Type)
	}

	config, err := s.buildLoggerOptions(id)
	if err != nil {
		return nil, errors.Join(errors.New("unable to build logger options"), err)
	}

	logger, err := provider.Logger(config)
	if err != nil {
		return nil, err
	}

	s.loggers[id] = logger

	return logger, nil
}

func (s *LoggerService) GetContextLogger(ctx context.Context, name string) (logging.Logger, error) {
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

func (s *LoggerService) GetLogger(id string) (logging.Logger, error) {
	logger, ok := s.loggers[id]
	if !ok {
		return nil, ErrUnknownLogger
	}
	return logger, nil
}

func (s *LoggerService) SetLevel(id, level string) error {
	logger, err := s.GetLogger(id)
	if err != nil {
		return errors.Join(errors.New("unable to set level"), err)
	}

	logger.SetLevel(level)
	return nil
}

// toPath combines the provided id and logsHome to create the log's path.
func (s *LoggerService) toPath(id string) string {
	return filepath.Join(s.config.LogsHome, fmt.Sprintf("%s.log", id))
}
