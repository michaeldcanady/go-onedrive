package di2

import "github.com/michaeldcanady/go-onedrive/internal/logging"

type LoggerService interface {
	CreateLogger(id string) (logging.Logger, error)
	SetLevel(id, level string) error
	GetLogger(id string) (logging.Logger, error)
	SetAllLevel(level string)
}
