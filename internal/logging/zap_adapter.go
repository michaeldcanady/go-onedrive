package logging

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Ensure ZapLogAdapter implements the Logger interface.
var _ Logger = (*ZapLogAdapter)(nil)

// ZapLogAdapter wraps a zap.Logger and exposes a dynamic log level.
type ZapLogAdapter struct {
	logger *zap.Logger
	level  *zap.AtomicLevel
}

// convertFieldsToZap converts a slice of custom Field types to zap.Fields.
func convertFieldsToZap(fields ...Field) ([]zap.Field, error) {
	zapFields := make([]zap.Field, len(fields))
	for index, field := range fields {
		switch field.fieldType {
		case FieldTypeString:
			zapFields[index] = zap.String(field.key, field.value.(string))
		case FieldTypeInt:
			zapFields[index] = zap.Int(field.key, field.value.(int))
		case FieldTypeAny:
			zapFields[index] = zap.Any(field.key, field.value)
		case FieldTypeBool:
			zapFields[index] = zap.Bool(field.key, field.value.(bool))
		default:
			return nil, fmt.Errorf("unknown field type: %v", field.fieldType)
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
	logger, err := cfg.Build()
	if err != nil {
		panic(err) // or return error if you prefer
	}

	return &ZapLogAdapter{
		logger: logger,
		level:  &atomicLevel,
	}
}

// NewZapLoggerAdapter wraps an existing zap.Logger.
// Note: this logger will NOT have a dynamically adjustable level
// unless you pass in a zap.Logger built with an AtomicLevel.
func NewZapLoggerAdapter(logger *zap.Logger) *ZapLogAdapter {
	return &ZapLogAdapter{
		logger: logger,
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

// Info logs a message at InfoLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (z *ZapLogAdapter) Debug(msg string, kv ...Field) {
	z.logger.Debug(msg, z.safeConvert(kv...)...)
}

// Error implements Logger.
func (z *ZapLogAdapter) Error(msg string, kv ...Field) {
	z.logger.Error(msg, z.safeConvert(kv...)...)
}

// Info logs a message at InfoLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (z *ZapLogAdapter) Info(msg string, kv ...Field) {
	z.logger.Info(msg, z.safeConvert(kv...)...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (z *ZapLogAdapter) Warn(msg string, kv ...Field) {
	z.logger.Warn(msg, z.safeConvert(kv...)...)
}

func (z *ZapLogAdapter) safeConvert(kv ...Field) []zap.Field {
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
