package infra

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Ensure ZapLogAdapter implements the domain.Logger interface.
var _ domain.Logger = (*ZapLogAdapter)(nil)

// ZapLogAdapter wraps a zap.Logger and exposes a dynamic log level.
type ZapLogAdapter struct {
	logger *zap.Logger
	level  *zap.AtomicLevel
	fields []domain.Field
}

// convertFieldsToZap converts a slice of custom Field types to zap.Fields.
func convertFieldsToZap(fields ...domain.Field) ([]zap.Field, error) {
	zapFields := make([]zap.Field, len(fields))
	for index, field := range fields {
		switch field.FieldType {
		case domain.FieldTypeString:
			zapFields[index] = zap.String(field.Key, field.Value.(string))
		case domain.FieldTypeInt:
			zapFields[index] = zap.Int(field.Key, field.Value.(int))
		case domain.FieldTypeAny:
			zapFields[index] = zap.Any(field.Key, field.Value)
		case domain.FieldTypeBool:
			zapFields[index] = zap.Bool(field.Key, field.Value.(bool))
		case domain.FieldTypeDuration:
			zapFields[index] = zap.Duration(field.Key, field.Value.(time.Duration))
		case domain.FieldTypeStrings:
			zapFields[index] = zap.Strings(field.Key, field.Value.([]string))
		case domain.FieldTypeError:
			zapFields[index] = zap.Error(field.Value.(error))
		case domain.FieldTypeTime:
			zapFields[index] = zap.Time(field.Key, field.Value.(time.Time))
		default:
			return nil, fmt.Errorf("unknown field type: %v", field.FieldType)
		}
	}
	return zapFields, nil
}

func convertLevelToZap(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zap.DebugLevel, nil
	case "info":
		return zap.InfoLevel, nil
	case "warn":
		return zap.WarnLevel, nil
	case "error":
		return zap.ErrorLevel, nil
	case "panic":
		return zap.PanicLevel, nil
	default:
		return -1, errors.New("invalid level")
	}
}

// NewZapLogger creates a new ZapLogAdapter with a dynamically adjustable log level.
func NewZapLogger(cfg zap.Config) *ZapLogAdapter {
	// Create an atomic level so it can be changed later
	atomicLevel := zap.NewAtomicLevel()

	// Override the config's level with the atomic level
	cfg.Level = atomicLevel

	// Build the logger
	l, err := cfg.Build()
	if err != nil {
		panic(err) // or return error if you prefer
	}

	return &ZapLogAdapter{
		logger: l,
		level:  &atomicLevel,
	}
}

// NewZapLoggerAdapter wraps an existing zap.Logger.
// Note: this logger will NOT have a dynamically adjustable level
// unless you pass in a zap.Logger built with an AtomicLevel.
func NewZapLoggerAdapter(l *zap.Logger) *ZapLogAdapter {
	return &ZapLogAdapter{
		logger: l,
		// level is zero-value (info) and not adjustable unless you expose it
	}
}

// SetLevel changes the logger's level at runtime.
func (z *ZapLogAdapter) SetLevel(level string) {
	if z.level == nil {
		z.logger.Sugar().Errorf("failed to set log level %s: %w", level, errors.New("atomicLevel empty"))
		return
	}

	zapLevel, err := convertLevelToZap(level)
	if err != nil {
		z.logger.Sugar().Errorf("failed to set log level %s: %w", level, err)
		return
	}
	z.level.SetLevel(zapLevel)
}

// Level returns the current log level.
func (z *ZapLogAdapter) Level() zapcore.Level {
	return z.level.Level()
}

// Debug logs a message at DebugLevel.
func (z *ZapLogAdapter) Debug(msg string, kv ...domain.Field) {
	z.logger.Debug(msg, z.safeConvert(kv...)...)
}

// Error logs a message at ErrorLevel.
func (z *ZapLogAdapter) Error(msg string, kv ...domain.Field) {
	z.logger.Error(msg, z.safeConvert(kv...)...)
}

// Info logs a message at InfoLevel.
func (z *ZapLogAdapter) Info(msg string, kv ...domain.Field) {
	z.logger.Info(msg, z.safeConvert(kv...)...)
}

// Warn logs a message at WarnLevel.
func (z *ZapLogAdapter) Warn(msg string, kv ...domain.Field) {
	z.logger.Warn(msg, z.safeConvert(kv...)...)
}

func (z *ZapLogAdapter) safeConvert(kv ...domain.Field) []zap.Field {
	if len(kv) == 0 {
		return nil
	}

	fields, err := convertFieldsToZap(kv...)
	if err != nil {
		z.logger.Sugar().Errorf("Failed to convert fields to zap: %v", err)
		return nil
	}

	return fields
}

func (z *ZapLogAdapter) WithContext(ctx context.Context) domain.Logger {
	ctxFields := domain.FromContextFields(ctx)
	if len(ctxFields) == 0 {
		return z
	}

	// Convert map → []Field
	newFields := make([]domain.Field, 0, len(ctxFields))
	for k, v := range ctxFields {
		newFields = append(newFields, domain.Any(k, v))
	}

	// Reuse the dedupe logic in With()
	return z.With(newFields...)
}

func (z *ZapLogAdapter) With(fields ...domain.Field) domain.Logger {
	if len(fields) == 0 {
		return z
	}

	// Deduplicate against existing fields
	newFields := make([]domain.Field, 0, len(fields))
	existing := make(map[string]bool)

	for _, f := range z.fields {
		existing[f.Key] = true
	}

	for _, f := range fields {
		if !existing[f.Key] {
			newFields = append(newFields, f)
		}
	}

	// Convert to zap fields
	zapFields, err := convertFieldsToZap(newFields...)
	if err != nil {
		z.logger.Sugar().Errorf("failed to convert fields in With: %v", err)
		return z
	}

	// Build new logger
	newZap := z.logger.With(zapFields...)

	// Return new adapter with merged fields
	return &ZapLogAdapter{
		logger: newZap,
		level:  z.level,
		fields: append(z.fields, newFields...),
	}
}
