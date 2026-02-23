package logging

import "context"

type NoopLogger struct{}

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

func (n *NoopLogger) Info(msg string, kv ...Field)  {}
func (n *NoopLogger) Warn(msg string, kv ...Field)  {}
func (n *NoopLogger) Error(msg string, kv ...Field) {}
func (n *NoopLogger) Debug(msg string, kv ...Field) {}

func (n *NoopLogger) SetLevel(level string) {}

func (n *NoopLogger) With(fields ...Field) Logger {
	return n
}

func (n *NoopLogger) WithContext(ctx context.Context) Logger {
	return n
}
