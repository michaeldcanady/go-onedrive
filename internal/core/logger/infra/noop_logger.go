package infra

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
)

var _ domain.Logger = (*NoopLogger)(nil)

type NoopLogger struct {
}

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

func (n *NoopLogger) Debug(msg string, kv ...domain.Field) {}
func (n *NoopLogger) Error(msg string, kv ...domain.Field) {}
func (n *NoopLogger) Info(msg string, kv ...domain.Field)  {}
func (n *NoopLogger) Warn(msg string, kv ...domain.Field)  {}
func (n *NoopLogger) SetLevel(level string)                {}
func (n *NoopLogger) With(fields ...domain.Field) domain.Logger {
	return n
}
func (n *NoopLogger) WithContext(ctx context.Context) domain.Logger {
	return n
}
