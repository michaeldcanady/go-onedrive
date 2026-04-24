package drive

import (
	"context"

	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
)

// DefaultResolver implements pkgfs.DriveResolver using the drive service.
type DefaultResolver struct {
	driveSvc   Service
	identityID string
}

// NewDefaultResolver initializes a new instance of the DefaultResolver.
func NewDefaultResolver(driveSvc Service, identityID string) pkgfs.DriveResolver {
	return &DefaultResolver{
		driveSvc:   driveSvc,
		identityID: identityID,
	}
}

// GetActiveDriveID retrieves the active drive ID.
// It defaults to the personal drive since explicit active drive selection was removed.
func (r *DefaultResolver) GetActiveDriveID(ctx context.Context) (string, error) {
	d, err := r.driveSvc.ResolvePersonalDrive(ctx, r.identityID)
	if err != nil {
		return "", err
	}
	return d.ID, nil
}
