package drive

import (
	"context"
)

type DriveService interface {
	ListDrives(ctx context.Context) ([]*Drive, error)
	ResolveDrive(ctx context.Context, ref string) (*Drive, error)
	ResolvePersonalDrive(ctx context.Context) (*Drive, error)
}
