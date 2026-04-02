package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapService manages a collection of Zap-backed Loggers.
type ZapService struct {
	// mu protects concurrent access to the loggers map.
	mu sync.RWMutex
	// loggers maps logger IDs to their respective ZapLogger instances.
	loggers map[string]*ZapLogger
	// root is the base logger used for creating named sub-loggers.
	root *zap.Logger
	// level is the global atomic level used by the root and sub-loggers.
	level zap.AtomicLevel
	// once ensures that the root logger is initialized only once.
	once sync.Once
}

// NewZapService initializes a new instance of the ZapService.
func NewZapService() *ZapService {
	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	return &ZapService{
		loggers: make(map[string]*ZapLogger),
		level:   level,
	}
}

// initRoot initializes the root zap logger if it hasn't been already.
// It configures the logger to output to a file by default.
func (s *ZapService) initRoot() error {
	var err error
	s.once.Do(func() {
		// Configure Zap logger to output to a file by default.
		// NOTE: This configuration does not automatically create the './logs' directory.
		// The environment must ensure this directory exists before the application starts.
		cfg := zap.Config{
			Encoding: "json", // Use JSON for structured logging
			Level:    s.level,
			//OutputPaths:      []string{"./logs/app.log"}, // Log to a file
			//ErrorOutputPaths: []string{"./logs/app.log"}, // Log errors to the same file
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey:    "message",
				LevelKey:      "level",
				TimeKey:       "time",
				CallerKey:     "caller",
				StacktraceKey: "stacktrace",
				EncodeTime:    zapcore.ISO8601TimeEncoder,
				EncodeLevel:   zapcore.CapitalLevelEncoder,
				EncodeCaller:  zapcore.ShortCallerEncoder,
			},
		}

		s.root, err = cfg.Build()
	})
	return err
}

// CreateLogger creates or returns an existing Logger for the given identifier.
func (s *ZapService) CreateLogger(id string) (Logger, error) {
	if err := s.initRoot(); err != nil {
		return nil, fmt.Errorf("failed to initialize root logger: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if l, ok := s.loggers[id]; ok {
		return l, nil
	}

	named := s.root.Named(id)
	l := &ZapLogger{
		logger: named,
		level:  s.level,
	}
	s.loggers[id] = l
	return l, nil
}

// SetAllLevel updates the severity level for all registered loggers.
func (s *ZapService) SetAllLevel(level string) {
	l, err := zapcore.ParseLevel(level)
	if err != nil {
		return
	}

	s.level.SetLevel(l)
}
