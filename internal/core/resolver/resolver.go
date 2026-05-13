// Package resolver provides utilities for resolving CLI arguments like paths, identities, and drives into their
// corresponding domain entities.
package resolver

import (
	"context"
	"fmt"
	"path"

	"github.com/michaeldcanady/go-onedrive/internal/features/config"
	"github.com/michaeldcanady/go-onedrive/internal/features/drive"
	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
	"github.com/michaeldcanady/go-onedrive/internal/features/vfs"
)

// Service coordinates the resolution of raw CLI arguments into domain-specific objects.
type Service interface {
	// ResolvePath translates a user-provided path string into a clean, absolute representation.
	ResolvePath(ctx context.Context, path string) (string, error)

	// ResolveIdentity searches for an identity matching the provided query, which can be an ID or a display name.
	ResolveIdentity(ctx context.Context, query string) (*identity.Identity, error)

	// ResolveDrive searches for a drive matching the provided query, which can be an ID or a name.
	ResolveDrive(ctx context.Context, query string) (*drive.Drive, error)
}

type resolverService struct {
	vfs        vfs.VFS
	identities identity.Service
	drives     drive.Service
	config     config.Service
}

// NewResolverService returns a new [Service] initialized with the provided dependencies.
func NewResolverService(v vfs.VFS, is identity.Service, ds drive.Service, cs config.Service) Service {
	return &resolverService{
		vfs:        v,
		identities: is,
		drives:     ds,
		config:     cs,
	}
}

func (s *resolverService) ResolvePath(ctx context.Context, p string) (string, error) {
	if path.IsAbs(p) {
		return path.Clean(p), nil
	}

	cwd := "/"
	if val, err := s.config.Get(config.KeyCoreVFSCWD); err == nil && val != nil {
		cwd = fmt.Sprintf("%v", val)
	}

	return path.Join(cwd, p), nil
}

func (s *resolverService) ResolveIdentity(ctx context.Context, query string) (*identity.Identity, error) {
	return s.identities.FindIdentity(ctx, query)
}

func (s *resolverService) ResolveDrive(ctx context.Context, query string) (*drive.Drive, error) {
	return s.drives.FindDrive(ctx, query)
}
