package shared

import (
	"context"

	"github.com/google/uuid"
)

type correlationIDKey struct{}

// NewCorrelationID generates a new correlation id.
func NewCorrelationID() string {
	return uuid.NewString()
}

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
