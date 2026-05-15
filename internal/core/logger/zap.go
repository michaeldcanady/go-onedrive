package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapService struct {
	logger *zap.Logger
	atomic zap.AtomicLevel
}

// NewZapLogger returns a new [Service] implementation backed by [zap.Logger].
// It automatically creates the directories for the specified outputPaths if they do not exist.
func NewZapLogger(outputPaths []string, level string) (Service, error) {
	atomic := zap.NewAtomicLevel()
	if err := atomic.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}

	config := zap.NewProductionConfig()
	config.Level = atomic
	config.OutputPaths = outputPaths
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Ensure directories for log files exist
	for _, path := range outputPaths {
		if path == "stdout" || path == "stderr" {
			continue
		}
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory %s: %w", dir, err)
		}
	}

	l, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build zap logger: %w", err)
	}

	return &zapService{
		logger: l,
		atomic: atomic,
	}, nil
}

func (z *zapService) Debug(msg string, fields ...any) {
	z.logger.Debug(msg, z.toZapFields(fields...)...)
}

func (z *zapService) Info(msg string, fields ...any) {
	z.logger.Info(msg, z.toZapFields(fields...)...)
}

func (z *zapService) Warn(msg string, fields ...any) {
	z.logger.Warn(msg, z.toZapFields(fields...)...)
}

func (z *zapService) Error(msg string, fields ...any) {
	z.logger.Error(msg, z.toZapFields(fields...)...)
}

func (z *zapService) Fatal(msg string, fields ...any) {
	z.logger.Fatal(msg, z.toZapFields(fields...)...)
}

func (z *zapService) With(fields ...any) Service {
	return &zapService{
		logger: z.logger.With(z.toZapFields(fields...)...),
		atomic: z.atomic,
	}
}

func (z *zapService) Sync() error {
	return z.logger.Sync()
}

func (z *zapService) SetLevel(level string) error {
	return z.atomic.UnmarshalText([]byte(level))
}

func (z *zapService) GetLevel() string {
	return z.atomic.String()
}

func (z *zapService) toZapFields(fields ...any) []zap.Field {
	numFields := len(fields)
	if numFields == 0 {
		return nil
	}
	zapFields := make([]zap.Field, 0, numFields/2)
	for i := 0; i < numFields; i += 2 {
		if i+1 < numFields {
			key, ok := fields[i].(string)
			if !ok {
				key = fmt.Sprintf("arg%d", i)
			}
			zapFields = append(zapFields, zap.Any(key, fields[i+1]))
		}
	}
	return zapFields
}
