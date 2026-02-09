package logger

import "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"

type LoggerProvider interface {
	Logger(opts logging.LoggerOptions) (logging.Logger, error)
}
