package logging

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
)

var _ logger.Logger = (*NoopLogger)(nil)

type NoopLogger struct {
}

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

func (n *NoopLogger) Debug(msg string, kv ...logger.Field) {}
func (n *NoopLogger) Error(msg string, kv ...logger.Field) {}
func (n *NoopLogger) Info(msg string, kv ...logger.Field)  {}
func (n *NoopLogger) Warn(msg string, kv ...logger.Field)  {}
func (n *NoopLogger) SetLevel(level string)                {}
func (n *NoopLogger) With(fields ...logger.Field) logger.Logger {
	return n
}
func (n *NoopLogger) WithContext(ctx context.Context) logger.Logger {
	return n
}
