package util

import (
	"context"

	"github.com/google/uuid"
)

type correlationIDKey struct{}

func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationIDKey{}, id)
}

func CorrelationIDFromContext(ctx context.Context) string {
	if v := ctx.Value(correlationIDKey{}); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewCorrelationID() string {
	return uuid.NewString()
}
