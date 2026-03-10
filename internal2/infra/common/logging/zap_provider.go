package logging

import (
	"errors"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"go.uber.org/zap"
)

type ZapLoggerProvider struct {
}

func NewZapLoggerProvider() *ZapLoggerProvider {
	return &ZapLoggerProvider{}
}

func (lP *ZapLoggerProvider) Logger(opts logger.LoggerOptions) (logger.Logger, error) {
	cfg := zap.NewProductionConfig()

	switch opts.Level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return nil, fmt.Errorf("unknown logging level: %v", opts.Level)
	}

	switch opts.OutputDestination {
	case logger.OutputDestinationStandardOut:
		cfg.OutputPaths = []string{"stdout"}
	case logger.OutputDestinationStandardError:
		cfg.OutputPaths = []string{"stderr"}
	case logger.OutputDestinationFile:
		logPath := opts.FilePath
		if logPath == "" {
			return nil, errors.New("log path must be specified for file output destination")
		}
		cfg.OutputPaths = []string{logPath}
	default:
		return nil, fmt.Errorf("unknown output destination: %v", opts.OutputDestination)
	}

	return NewZapLogger(cfg), nil
}
