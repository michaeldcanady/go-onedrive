package di

import (
	"context"

	driveservice2 "github.com/michaeldcanady/go-onedrive/internal/app/drive_service2"
)

type DriveService interface {
	ListDrives(ctx context.Context) ([]*driveservice2.Drive, error)
	ResolveDrive(ctx context.Context, driveRef string) (*driveservice2.Drive, error)
	ResolvePersonalDrive(ctx context.Context) (*driveservice2.Drive, error)
}
