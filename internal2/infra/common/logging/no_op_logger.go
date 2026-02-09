package logging

import "context"

type NoopLogger struct{}

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

// SetLevel implements [Logger].
func (l *NoopLogger) SetLevel(_ string) {}

// With implements [Logger].
func (l *NoopLogger) With(_ ...Field) Logger { return l }

// WithContext implements [Logger].
func (l *NoopLogger) WithContext(_ context.Context) Logger { return l }

func (l *NoopLogger) Debug(_ string, _ ...Field) {}
func (l *NoopLogger) Info(_ string, _ ...Field)  {}
func (l *NoopLogger) Warn(_ string, _ ...Field)  {}
func (l *NoopLogger) Error(_ string, _ ...Field) {}
