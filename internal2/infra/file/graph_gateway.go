package file

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
)

// GraphGateway defines the interface for interacting with the Microsoft Graph API for file metadata.
type GraphGateway interface {
	GetByPath(ctx context.Context, driveID, path string, etag string) (*file.Metadata, error)
	ListByPath(ctx context.Context, driveID, path string, parentEtag string) ([]*file.Metadata, error)
	CreateByPath(ctx context.Context, driveID, parentPath string, request file.MetadataCreateRequest) (*file.Metadata, error)
	UpdateByPath(ctx context.Context, driveID, path string, request file.MetadataUpdateRequest) (*file.Metadata, error)
	DeleteByPath(ctx context.Context, driveID, path string) error
}
