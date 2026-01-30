package logging

import (
	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
)

var _ logger.LoggerProvider = (*LoggerProvider)(nil)

type LoggerProvider struct {
}

func NewLoggerProvider() *LoggerProvider {
	return &LoggerProvider{}
}

// Logger implements [logger.LoggerProvider].
func (l *LoggerProvider) Logger(logLevel string, logsHome string) (logging.Logger, error) {
	return logging.NewZapProvider(logLevel, logsHome)
}
