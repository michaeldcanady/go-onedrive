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
}

// NewZapService initializes a new instance of the ZapService.
func NewZapService() *ZapService {
	return &ZapService{
		loggers: make(map[string]*ZapLogger),
	}
}

// CreateLogger creates or returns an existing Logger for the given identifier.
func (s *ZapService) CreateLogger(id string) (Logger, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if l, ok := s.loggers[id]; ok {
		return l, nil
	}

	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	cfg := zap.NewProductionConfig()
	cfg.Level = level

	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build zap logger: %w", err)
	}

	l := &ZapLogger{
		logger: zapLogger,
		level:  level,
	}
	s.loggers[id] = l
	return l, nil
}

// SetAllLevel updates the severity level for all registered loggers.
func (s *ZapService) SetAllLevel(level string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	l, err := zapcore.ParseLevel(level)
	if err != nil {
		return
	}

	for _, logger := range s.loggers {
		logger.level.SetLevel(l)
	}
}
