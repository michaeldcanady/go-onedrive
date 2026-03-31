package logger

import (
	"context"
)

// Logger is an interface for logging messages with various severity levels.
// It provides a consistent logging API across the application.
type Logger interface {
	// Info logs a message at Info level.
	Info(msg string, kv ...Field)
	// Warn logs a message at Warning level.
	Warn(msg string, kv ...Field)
	// Error logs a message at Error level.
	Error(msg string, kv ...Field)
	// Debug logs a message at Debug level.
	Debug(msg string, kv ...Field)
	// SetLevel dynamically changes the logging level for this logger instance.
	SetLevel(level string)
	// With returns a new logger with the provided fields pre-populated.
	With(fields ...Field) Logger
	// WithContext returns a new logger that is aware of the provided context.
	WithContext(ctx context.Context) Logger
}
