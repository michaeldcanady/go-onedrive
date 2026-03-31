package middleware

import (
	"net/http"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/feature/logger"
	nethttp "github.com/microsoft/kiota-http-go"
)

// KiotaLoggingMiddleware implements the Kiota Middleware interface to provide
// detailed request and response logging.
type KiotaLoggingMiddleware struct {
	log logger.Logger
}

// NewKiotaLoggingMiddleware initializes a new instance of KiotaLoggingMiddleware.
func NewKiotaLoggingMiddleware(log logger.Logger) *KiotaLoggingMiddleware {
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
			logger.Duration("duration", duration),
		)
		return resp, err
	}

	log.Debug("inbound response",
		logger.String("method", req.Method),
		logger.String("url", req.URL.String()),
		logger.Int("status", resp.StatusCode),
		logger.Duration("duration", duration),
	)

	if resp.StatusCode >= 400 {
		log.Warn("request returned error status",
			logger.Int("status", resp.StatusCode),
			logger.String("url", req.URL.String()),
		)
	}

	return resp, nil
}
