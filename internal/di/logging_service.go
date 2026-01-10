package di

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/logging"
)

type LoggingService interface {
	Info(ctx context.Context, msg string, kv ...logging.Field) error
	Warn(ctx context.Context, msg string, kv ...logging.Field) error
	Error(ctx context.Context, msg string, kv ...logging.Field) error
	Debug(ctx context.Context, msg string, kv ...logging.Field) error
}
