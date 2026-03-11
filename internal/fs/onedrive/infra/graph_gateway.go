package infra

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/fs/shared/domain"
)

// GraphGateway defines the interface for interacting with the Microsoft Graph API for file metadata.
type GraphGateway interface {
	GetByPath(ctx context.Context, driveID, path string, etag string) (*domain.Metadata, error)
	ListByPath(ctx context.Context, driveID, path string, parentEtag string) ([]*domain.Metadata, error)
	CreateByPath(ctx context.Context, driveID, parentPath string, request domain.MetadataCreateRequest) (*domain.Metadata, error)
	UpdateByPath(ctx context.Context, driveID, path string, request domain.MetadataUpdateRequest) (*domain.Metadata, error)
	DeleteByPath(ctx context.Context, driveID, path string) error
}
