package infra

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/fs/shared/domain"
)

// GraphContentsGateway defines the interface for interacting with Microsoft Graph for file content.
type GraphContentsGateway interface {
	Download(ctx context.Context, driveID, path string, etag string) (fresh []byte, ctag string, err error)
	Upload(ctx context.Context, driveID, path string, data []byte, ifMatch string) (*domain.Metadata, string, error)
}
