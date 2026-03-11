package app

import (
	"context"

	domaindrive "github.com/michaeldcanady/go-onedrive/internal/drive/domain"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
)

var _ domaindrive.DriveResolver = (*DriverResolverAdapter)(nil)

type DriverResolverAdapter struct {
	stateSvc     domainstate.Service
	driveService domaindrive.DriveService
}

// CurrentDriveID implements [domaindrive.DriveResolver].
func (d *DriverResolverAdapter) CurrentDriveID(ctx context.Context) (string, error) {
	id, err := d.stateSvc.Get(domainstate.KeyDrive)
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

func NewDriverResolverAdapter(stateSvc domainstate.Service, driveService domaindrive.DriveService) *DriverResolverAdapter {
	return &DriverResolverAdapter{
		stateSvc:     stateSvc,
		driveService: driveService,
	}
}
