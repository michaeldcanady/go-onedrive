package file

import (
	"context"
	"io"

	infrafile "github.com/michaeldcanady/go-onedrive/internal2/infra/file"
)

type FileService interface {
	ResolveItem(ctx context.Context, driveID, path string) (*infrafile.DriveItem, error)
	ListChildren(ctx context.Context, driveID, path string) ([]*infrafile.DriveItem, error)
	GetFileContents(ctx context.Context, driveID, path string) ([]byte, error)
	WriteFile(ctx context.Context, driveID, path string, reader io.Reader) (*infrafile.DriveItem, error)
	// later: Upload, Download, Delete, Move
}
