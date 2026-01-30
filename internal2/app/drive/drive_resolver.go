package drive

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
)

var _ drive.DriveResolver = (*DriverResolverAdapter)(nil)

type DriverResolverAdapter struct {
	driveService *driveService
}

func NewDriverResolverAdapter(driveService *driveService) *DriverResolverAdapter {
	return &DriverResolverAdapter{
		driveService: driveService,
	}
}

// CurrentDriveID implements [drive.DriveResolver].
func (d *DriverResolverAdapter) CurrentDriveID(ctx context.Context) (string, error) {
	drive, err := d.driveService.ResolvePersonalDrive(ctx)
	if err != nil {
		return "", err
	}

	return drive.ID, nil
}
