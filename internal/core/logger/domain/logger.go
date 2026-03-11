package domain

import "context"

// Logger is an interface for logging messages with various severity levels.
type Logger interface {
	// Info logs a message at InfoLevel.
	Info(msg string, kv ...Field)
	// Warn logs a message at WarnLevel.
	Warn(msg string, kv ...Field)
	// Error logs a message at ErrorLevel.
	Error(msg string, kv ...Field)
	// Debug logs a message at DebugLevel.
	Debug(msg string, kv ...Field)

	// SetLevel Change the level of the log.
	SetLevel(level string)

	With(fields ...Field) Logger
	WithContext(ctx context.Context) Logger
}
