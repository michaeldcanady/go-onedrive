package state

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/state"
)

var _ drive.DriveResolver = (*DriverResolverAdapter)(nil)

type DriverResolverAdapter struct {
	stateSvc state.Service
}

// CurrentDriveID implements [drive.DriveResolver].
func (d *DriverResolverAdapter) CurrentDriveID(ctx context.Context) (string, error) {
	return d.stateSvc.GetCurrentDrive()
}

func NewDriverResolverAdapter(stateSvc state.Service) *DriverResolverAdapter {
	return &DriverResolverAdapter{
		stateSvc: stateSvc,
	}
}
