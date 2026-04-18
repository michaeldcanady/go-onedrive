package alias

import (
	"context"
)

// Repository defines the persistence interface for drive aliases.
type Repository interface {
	Get(ctx context.Context, name string) (string, error)
	Set(ctx context.Context, name, driveID string) error
	Delete(ctx context.Context, name string) error
	List(ctx context.Context) (map[string]string, error)
}
