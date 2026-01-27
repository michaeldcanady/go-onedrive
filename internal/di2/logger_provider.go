package di2

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/logging"
	"go.uber.org/zap"
)

type loggerProvider struct{}

func (p *loggerProvider) Logger(logLevel string, logsHome string) (logging.Logger, error) {
	return p.createZapLogger(logLevel, logsHome)
}

func (p *loggerProvider) createZapLogger(logLevel, logsHome string) (logging.Logger, error) {
	cfg := zap.NewProductionConfig()

	switch logLevel {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return nil, fmt.Errorf("unknown logging level: %s", logLevel)
	}

	// Logs to stdout
	//cfg.OutputPaths = append(cfg.OutputPaths, logsHome)
	// only log to file
	cfg.OutputPaths = []string{logsHome}

	return logging.NewZapLogger(cfg), nil
}
