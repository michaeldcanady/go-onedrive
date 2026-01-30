package logger

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
)

type LoggerService interface {
	CreateLogger(id string) (logging.Logger, error)
	GetContextLogger(ctx context.Context, name string) (logging.Logger, error)
	GetLogger(id string) (logging.Logger, error)
	SetAllLevel(level string)
	SetLevel(id string, level string) error
}
