package file

import (
	"context"
	"io"
)

type FileContentsRepository interface {
	Download(ctx context.Context, driveID string, path string, opts DownloadOptions) (io.ReadCloser, error)
	Upload(ctx context.Context, driveID string, path string, r io.Reader, opts UploadOptions) (*Metadata, error)
}
