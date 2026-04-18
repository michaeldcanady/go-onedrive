package drive

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/state"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
)

// DefaultResolver implements pkgfs.DriveResolver using the internal state service.
type DefaultResolver struct {
	state      state.Service
	identityID string
}

// NewDefaultResolver initializes a new instance of the DefaultResolver.
func NewDefaultResolver(state state.Service, identityID string) pkgfs.DriveResolver {
	return &DefaultResolver{
		state:      state,
		identityID: identityID,
	}
}

// GetActiveDriveID retrieves the active drive ID from state, scoped to an identity if provided.
func (r *DefaultResolver) GetActiveDriveID(ctx context.Context) (string, error) {
	if r.identityID != "" {
		return r.state.GetScoped("tokens/microsoft", r.identityID+"/active_drive")
	}
	return r.state.Get(state.KeyDrive)
}
