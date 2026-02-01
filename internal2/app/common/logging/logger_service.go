package logging

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"

	domainlogger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
)

type LoggerService struct {
	loggers  map[string]logging.Logger
	logLevel string
	logsHome string
	factory  domainlogger.LoggerProvider
}

func NewLoggerService(logLevel string, logsHome string, factory domainlogger.LoggerProvider) (*LoggerService, error) {

	if err := os.MkdirAll(logsHome, os.ModePerm); err != nil {
		return nil, err
	}

	return &LoggerService{
		loggers:  map[string]logging.Logger{},
		logLevel: logLevel,
		logsHome: logsHome,
		factory:  factory,
	}, nil
}

// createLogger creates a new logger using the current factory.
func (s *LoggerService) CreateLogger(id string) (logging.Logger, error) {
	logger, err := s.factory.Logger(s.logLevel, s.toPath(id))
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
	s.logLevel = level
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
	return filepath.Join(s.logsHome, fmt.Sprintf("%s.log", id))
}
