package drive

import (
	"context"
)

type DriveGateway interface {
	ListDrives(context.Context) ([]*Drive, error)
	GetPersonalDrive(context.Context) (*Drive, error)
}
