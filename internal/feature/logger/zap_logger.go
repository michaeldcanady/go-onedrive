package logger

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/feature/shared"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger is a concrete implementation of the Logger interface using the Zap logging library.
type ZapLogger struct {
	// logger is the underlying Zap logger instance.
	logger *zap.Logger
	// level is an atomic level that allows for dynamic level changes.
	level zap.AtomicLevel
}

// Info logs a message at the Info severity level.
func (z *ZapLogger) Info(msg string, kv ...Field) {
	z.logger.Info(msg, z.toZapFields(kv...)...)
}

// Warn logs a message at the Warning severity level.
func (z *ZapLogger) Warn(msg string, kv ...Field) {
	z.logger.Warn(msg, z.toZapFields(kv...)...)
}

// Error logs a message at the Error severity level.
func (z *ZapLogger) Error(msg string, kv ...Field) {
	z.logger.Error(msg, z.toZapFields(kv...)...)
}

// Debug logs a message at the Debug severity level.
func (z *ZapLogger) Debug(msg string, kv ...Field) {
	z.logger.Debug(msg, z.toZapFields(kv...)...)
}

// SetLevel dynamically updates the logger's severity level.
func (z *ZapLogger) SetLevel(level string) {
	l, err := zapcore.ParseLevel(level)
	if err == nil {
		z.level.SetLevel(l)
	}
}

// With returns a new ZapLogger that includes the provided fields.
func (z *ZapLogger) With(fields ...Field) Logger {
	return &ZapLogger{
		logger: z.logger.With(z.toZapFields(fields...)...),
		level:  z.level,
	}
}

// WithContext returns a logger that is contextualized for the given context.
func (z *ZapLogger) WithContext(ctx context.Context) Logger {
	cid := shared.CorrelationIDFromContext(ctx)
	if cid == "" {
		return z
	}

	return &ZapLogger{
		logger: z.logger.With(zap.String("correlation_id", cid)),
		level:  z.level,
	}
}

// toZapFields transforms internal Field types into Zap-compatible fields.
func (z *ZapLogger) toZapFields(fields ...Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		switch f.FieldType {
		case FieldTypeString:
			zapFields[i] = zap.String(f.Key, f.Value.(string))
		case FieldTypeInt:
			zapFields[i] = zap.Int(f.Key, f.Value.(int))
		case FieldTypeTime:
			zapFields[i] = zap.Time(f.Key, f.Value.(time.Time))
		case FieldTypeDuration:
			zapFields[i] = zap.Duration(f.Key, f.Value.(time.Duration))
		case FieldTypeError:
			zapFields[i] = zap.Error(f.Value.(error))
		case FieldTypeBool:
			zapFields[i] = zap.Bool(f.Key, f.Value.(bool))
		}
	}
	return zapFields
}
