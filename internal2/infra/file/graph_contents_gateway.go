package file

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
)

// GraphContentsGateway defines the interface for interacting with Microsoft Graph for file content.
type GraphContentsGateway interface {
	Download(ctx context.Context, driveID, path string, etag string) (fresh []byte, ctag string, err error)
	Upload(ctx context.Context, driveID, path string, data []byte, ifMatch string) (*file.Metadata, string, error)
}
