package di

import (
	"context"

	driveservice "github.com/michaeldcanady/go-onedrive/internal/app/file_service"
)

type FileSystemService interface {
	ResolveItem(ctx context.Context, driveID, path string) (*driveservice.DriveItem, error)
	ListChildren(ctx context.Context, driveID, path string) ([]*driveservice.DriveItem, error)
}
