package state

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/state"
)

var _ drive.DriveResolver = (*DriverResolverAdapter)(nil)

type DriverResolverAdapter struct {
	stateSvc     state.Service
	driveService drive.DriveService
}

// CurrentDriveID implements [drive.DriveResolver].
func (d *DriverResolverAdapter) CurrentDriveID(ctx context.Context) (string, error) {
	id, err := d.stateSvc.GetCurrentDrive()
	if err != nil {
		return "", err
	}

	if id != "" {
		return id, nil
	}

	// Fallback to personal drive
	pDrive, err := d.driveService.ResolvePersonalDrive(ctx)
	if err != nil {
		return "", err
	}

	return pDrive.ID, nil
}

func NewDriverResolverAdapter(stateSvc state.Service, driveService drive.DriveService) *DriverResolverAdapter {
	return &DriverResolverAdapter{
		stateSvc:     stateSvc,
		driveService: driveService,
	}
}
