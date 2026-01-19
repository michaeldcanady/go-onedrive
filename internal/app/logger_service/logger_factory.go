package loggerservice

import "github.com/michaeldcanady/go-onedrive/internal/logging"

// loggerFactory represents a factory for creating a new logger.
type loggerFactory func(level, logPath string) (logging.Logger, error)
