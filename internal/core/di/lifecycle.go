package di

import (
	"context"
)

// Initializer is an interface for services that require explicit initialization.
type Initializer interface {
	Init(ctx context.Context) error
}

// Shutdowner is an interface for services that require explicit cleanup.
type Shutdowner interface {
	Shutdown(ctx context.Context) error
}
