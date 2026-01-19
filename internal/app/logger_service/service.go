package loggerservice

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"

	"github.com/michaeldcanady/go-onedrive/internal/logging"
)

type Service struct {
	loggers  map[string]logging.Logger
	logLevel string
	logsHome string
	factory  loggerFactory
}

func New(logLevel string, logsHome string) (*Service, error) {

	if err := os.MkdirAll(logsHome, os.ModePerm); err != nil {
		return nil, err
	}

	return &Service{
		loggers:  map[string]logging.Logger{},
		logLevel: logLevel,
		logsHome: logsHome,
		factory:  newZapLogger,
	}, nil
}

// createLogger creates a new logger using the current factory.
func (s *Service) CreateLogger(id string) (logging.Logger, error) {
	logger, err := s.factory(s.logLevel, s.toPath(id))
	if err != nil {
		return nil, err
	}

	s.loggers[id] = logger

	return logger, nil
}

func (s *Service) SetAllLevel(level string) {
	maps.Keys(s.loggers)(func(key string) bool {
		s.SetLevel(key, level)
		return true
	})
}

func (s *Service) GetLogger(id string) (logging.Logger, error) {
	logger, ok := s.loggers[id]
	if !ok {
		return nil, errors.New("unknown logger id")
	}
	return logger, nil
}

func (s *Service) SetLevel(id, level string) error {
	logger, err := s.GetLogger(id)
	if err != nil {
		return errors.Join(errors.New("unable to set level"), err)
	}

	logger.SetLevel(level)
	return nil
}

// toPath combines the provided id and logsHome to create the log's path.
func (s *Service) toPath(id string) string {
	return filepath.Join(s.logsHome, fmt.Sprintf("%s.log", id))
}
