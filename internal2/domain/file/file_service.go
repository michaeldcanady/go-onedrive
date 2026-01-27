package file

import (
	"context"

	infrafile "github.com/michaeldcanady/go-onedrive/internal2/infra/file"
)

type FileService interface {
	ResolveItem(ctx context.Context, driveID, path string) (*infrafile.DriveItem, error)
	ListChildren(ctx context.Context, driveID, path string) ([]*infrafile.DriveItem, error)
	// later: Upload, Download, Delete, Move
}
