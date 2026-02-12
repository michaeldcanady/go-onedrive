package logging

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type ZapLoggerProvider struct {
}

func NewZapLoggerProvider() *ZapLoggerProvider {
	return &ZapLoggerProvider{}
}

func (lP *ZapLoggerProvider) Logger(opts LoggerOptions) (Logger, error) {
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

	switch opts.outputDestination {
	case OutputDestinationStandardOut:
		cfg.OutputPaths = []string{"stdout"}
	case OutputDestinationStandardError:
		cfg.OutputPaths = []string{"stderr"}
	case OutputDestinationFile:
		logPath := opts.logPath
		if logPath == "" {
			return nil, errors.New("log path must be specified for file output destination")
		}
		cfg.OutputPaths = []string{logPath}
	default:
		return nil, fmt.Errorf("unknown output destination: %v", opts.outputDestination)
	}

	return NewZapLogger(cfg), nil
}
