package util

import (
	applogging "github.com/michaeldcanady/go-onedrive/internal2/app/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
)

// EnsureLogger retrieves or creates the CLI logger.
func EnsureLogger(c di.Container, name string) (infralogging.Logger, error) {
	logger, err := c.Logger().GetLogger(name)
	if err == applogging.ErrUnknownLogger {
		return c.Logger().CreateLogger(name)
	}
	return logger, err
}
