package logger

import (
	"context"
)

type LoggerService interface {
	CreateLogger(id string) (Logger, error)
	GetContextLogger(ctx context.Context, name string) (Logger, error)
	GetLogger(id string) (Logger, error)
	SetAllLevel(level string)
	SetLevel(id string, level string) error
	RegisterProvider(ype Type, factory LoggerProvider)
}
