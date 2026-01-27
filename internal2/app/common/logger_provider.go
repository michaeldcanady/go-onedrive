package common

import "github.com/michaeldcanady/go-onedrive/internal/logging"

type LoggerProvider interface {
	Logger(logLevel string, logsHome string) (logging.Logger, error)
}
