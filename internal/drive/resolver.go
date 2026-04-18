package drive

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/profile"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
)

// DefaultResolver implements pkgfs.DriveResolver using the profile service.
type DefaultResolver struct {
	profileSvc profile.Service
	identityID string
}

// NewDefaultResolver initializes a new instance of the DefaultResolver.
func NewDefaultResolver(profileSvc profile.Service, identityID string) pkgfs.DriveResolver {
	return &DefaultResolver{
		profileSvc: profileSvc,
		identityID: identityID,
	}
}

// GetActiveDriveID retrieves the active drive ID for the current profile.
func (r *DefaultResolver) GetActiveDriveID(ctx context.Context) (string, error) {
	p, err := r.profileSvc.GetActive(ctx)
	if err != nil {
		return "", err
	}
	return p.ActiveDriveID, nil
}
