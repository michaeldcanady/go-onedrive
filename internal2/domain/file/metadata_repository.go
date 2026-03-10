package file

import (
	"context"
)

type MetadataRepository interface {
	GetByPath(ctx context.Context, driveID string, path string, opts MetadataGetOptions) (*Metadata, error)
	ListByPath(ctx context.Context, driveID string, path string, opts MetadataListOptions) ([]*Metadata, error)
	CreateByPath(ctx context.Context, driveID, parentPath string, body MetadataCreateRequest, opts MetadataCreateOptions) (*Metadata, error)
	UpdateByPath(ctx context.Context, driveID, path string, body MetadataUpdateRequest, opts MetadataUpdateOptions) (*Metadata, error)
	DeleteByPath(ctx context.Context, driveID, path string, opts MetadataDeleteOptions) error
}
