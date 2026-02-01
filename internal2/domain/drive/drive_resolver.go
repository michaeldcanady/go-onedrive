package drive

import "context"

type DriveResolver interface {
	CurrentDriveID(ctx context.Context) (string, error)
}
