// Package logger provides a structured logging abstraction and utilities for context-aware logging.
package logger

import "context"

// Service coordinates structured logging throughout the application.
// It is safe for concurrent use by multiple goroutines.
type Service interface {
	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)
	Fatal(msg string, fields ...any)

	// With returns a new [Service] with the provided fields added to its context.
	With(fields ...any) Service

	// Sync flushes any buffered log entries.
	Sync() error

	// SetLevel dynamically changes the logging level at runtime.
	SetLevel(level string) error

	// GetLevel returns the current logging level.
	GetLevel() string
}

// Key for request ID in context
type requestIDKey struct{}

// WithContext returns a [Service] enriched with metadata from the context, such as request IDs.
// It does NOT retrieve the logger from the context itself.
func WithContext(l Service, ctx context.Context) Service {
	if l == nil {
		return nil
	}
	if rid := GetRequestID(ctx); rid != "" {
		return l.With("request_id", rid)
	}
	return l
}

// WithRequestID returns a new context containing the provided request ID.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// GetRequestID returns the request ID from the context, or an empty string if none exists.
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		return id
	}
	return ""
}
