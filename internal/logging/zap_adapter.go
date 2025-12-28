package logging

import (
	"fmt"

	"go.uber.org/zap"
)

// Ensure ZapLogAdapter implements the Logger interface.
var _ Logger = (*ZapLogAdapter)(nil)

// ZapLogAdapter is an adapter that implements the Logger interface using a zap.Logger.
type ZapLogAdapter struct {
	logger *zap.Logger
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
		default:
			return nil, fmt.Errorf("unknown field type: %v", field.fieldType)
		}
	}
	return zapFields, nil
}

// NewZapLoggerAdapter creates a new ZapLogAdapter that wraps a zap.Logger.
func NewZapLoggerAdapter(logger *zap.Logger) *ZapLogAdapter {
	return &ZapLogAdapter{
		logger: logger,
	}
}

// Info logs a message at InfoLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (z *ZapLogAdapter) Debug(msg string, kv ...Field) {
	if len(kv) == 0 {
		z.logger.Debug(msg)
		return
	}
	zapFields, err := convertFieldsToZap(kv...)
	if err != nil {
		z.logger.Sugar().Errorf("Failed to convert fields to zap: %v", err)
		return
	}
	z.logger.Debug(msg, zapFields...)
}

// Error implements Logger.
func (z *ZapLogAdapter) Error(msg string, kv ...Field) {
	if len(kv) == 0 {
		z.logger.Error(msg)
		return
	}
	zapFields, err := convertFieldsToZap(kv...)
	if err != nil {
		z.logger.Sugar().Errorf("Failed to convert fields to zap: %v", err)
		return
	}
	z.logger.Error(msg, zapFields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (z *ZapLogAdapter) Info(msg string, kv ...Field) {
	if len(kv) == 0 {
		z.logger.Info(msg)
		return
	}
	zapFields, err := convertFieldsToZap(kv...)
	if err != nil {
		z.logger.Sugar().Errorf("Failed to convert fields to zap: %v", err)
		return
	}
	z.logger.Info(msg, zapFields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (z *ZapLogAdapter) Warn(msg string, kv ...Field) {
	if len(kv) == 0 {
		z.logger.Warn(msg)
		return
	}
	zapFields, err := convertFieldsToZap(kv...)
	if err != nil {
		z.logger.Sugar().Errorf("Failed to convert fields to zap: %v", err)
		return
	}
	z.logger.Warn(msg, zapFields...)
}
