package alias

import (
	"context"
)

// Service provides methods to manage drive aliases.
type Service interface {
	GetDriveIDByAlias(ctx context.Context, name string) (string, error)
	GetAliasByDriveID(ctx context.Context, driveID string) (string, error)
	SetAlias(ctx context.Context, name, driveID string) error
	DeleteAlias(ctx context.Context, name string) error
	ListAliases(ctx context.Context) (map[string]string, error)
	Close() error
}
