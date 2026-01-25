package loggerservice

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/logging"
	"go.uber.org/zap"
)

func newZapLogger(logLevel string, logsHome string) (logging.Logger, error) {
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

	cfg.OutputPaths = append(cfg.OutputPaths, logsHome)

	return logging.NewZapLogger(cfg), nil
}
