package file

import (
	"context"
	"io"
)

type FileService interface {
	ResolveItem(ctx context.Context, driveID, path string) (*DriveItem, error)
	ListChildren(ctx context.Context, driveID, path string) ([]*DriveItem, error)
	GetFileContents(ctx context.Context, driveID, path string) ([]byte, error)
	WriteFile(ctx context.Context, driveID, path string, reader io.Reader) (*DriveItem, error)
	// later: Upload, Download, Delete, Move
}
