package di2

import (
	"context"

	fileservice "github.com/michaeldcanady/go-onedrive/internal/app/file_service"
)

type FileService interface {
	ResolveItem(ctx context.Context, driveID, path string) (*fileservice.DriveItem, error)
	ListChildren(ctx context.Context, driveID, path string) ([]*fileservice.DriveItem, error)
	// later: Upload, Download, Delete, Move
}
