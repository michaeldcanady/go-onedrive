package util

import (
	"context"

	applogging "github.com/michaeldcanady/go-onedrive/internal2/app/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
)

// EnsureLogger retrieves or creates the CLI logger.
func EnsureLogger(ctx context.Context, c di.Container, name string) (infralogging.Logger, error) {
	logger, err := c.Logger().GetLogger(name)
	if err == applogging.ErrUnknownLogger {
		logger, err = c.Logger().CreateLogger(name)
	}

	if cid := CorrelationIDFromContext(ctx); cid != "" {
		if logger != nil {
			logger = logger.With(logging.String("correlation_id", cid))
		}
	}

	return logger, err
}
