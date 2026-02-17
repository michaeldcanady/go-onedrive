package file

import (
	"context"
	"io"
)

type FileContentsRepository interface {
	Download(ctx context.Context, id string, opts DownloadOptions) (io.ReadCloser, string, error)
	Upload(ctx context.Context, id string, r io.Reader, opts UploadOptions) (string, error)
}
