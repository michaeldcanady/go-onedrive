package logger

import "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"

type LoggerProvider interface {
	Logger(logLevel string, logsHome string) (logging.Logger, error)
}
