package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	nethttp "github.com/microsoft/kiota-http-go"
)

// Logger defines the interface required for logging within the middleware.
type Logger interface {
	Debug(msg string, fields ...logger.Field)
	Error(msg string, fields ...logger.Field)
	WithContext(ctx context.Context) logger.Logger
}

// KiotaLoggingMiddleware implements the Kiota Middleware interface to provide
// detailed request and response logging.
type KiotaLoggingMiddleware struct {
	log Logger
}

// NewKiotaLoggingMiddleware initializes a new instance of KiotaLoggingMiddleware.
func NewKiotaLoggingMiddleware(log Logger) *KiotaLoggingMiddleware {
	return &KiotaLoggingMiddleware{
		log: log,
	}
}

// Intercept captures the request and response to log diagnostic information.
func (m *KiotaLoggingMiddleware) Intercept(pipeline nethttp.Pipeline, middlewareIndex int, req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	log := m.log.WithContext(ctx)

	log.Debug("outbound request",
		logger.String("method", req.Method),
		logger.String("url", req.URL.String()),
	)

	start := time.Now()
	resp, err := pipeline.Next(req, middlewareIndex)
	duration := time.Since(start)

	if err != nil {
		log.Error("request failed",
			logger.String("method", req.Method),
			logger.String("url", req.URL.String()),
			logger.Error(err),
		)
		return nil, err
	}

	log.Debug("inbound response",
		logger.String("method", req.Method),
		logger.String("url", req.URL.String()),
		logger.Int("status", resp.StatusCode),
		logger.Duration("duration", duration),
	)

	return resp, nil
}
