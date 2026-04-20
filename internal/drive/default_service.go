package drive

import (
	"context"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// DefaultService provides the default implementation of the drive service.
type DefaultService struct {
	vfs *fs.VFS
	log logger.Logger
}

// NewDefaultService initializes a new instance of the DefaultService.
func NewDefaultService(vfs *fs.VFS, l logger.Logger) *DefaultService {
	return &DefaultService{
		vfs: vfs,
		log: l,
	}
}

// ListDrives retrieves all accessible OneDrive drives.
func (s *DefaultService) ListDrives(ctx context.Context, identityID string) ([]Drive, error) {
	s.log.Debug("listing drives", logger.String("identity", identityID))

	// Assuming '/personal' is the mount path for now.
	drives, err := s.vfs.ListDrives(ctx, "/personal")
	if err != nil {
		s.log.Error("failed to list drives", logger.Error(err))
		return nil, fmt.Errorf("failed to list drives: %w", err)
	}

	var out []Drive
	for _, d := range drives {
		out = append(out, Drive{
			ID:       d.ID,
			Name:     d.Name,
			Type:     d.Type,
			Owner:    d.Owner,
			ReadOnly: d.ReadOnly,
		})
	}
	return out, nil
}

// ResolveDrive identifies a drive by its ID or name.
func (s *DefaultService) ResolveDrive(ctx context.Context, driveRef string, identityID string) (Drive, error) {
	s.log.Debug("resolving drive", logger.String("ref", driveRef), logger.String("identity", identityID))

	drives, err := s.ListDrives(ctx, identityID)
	if err != nil {
		return Drive{}, err
	}

	for _, d := range drives {
		if strings.EqualFold(d.ID, driveRef) || strings.EqualFold(d.Name, driveRef) {
			return d, nil
		}
	}

	return Drive{}, fmt.Errorf("drive %s not found", driveRef)
}

// ResolvePersonalDrive retrieves the user's primary personal drive.
func (s *DefaultService) ResolvePersonalDrive(ctx context.Context, identityID string) (Drive, error) {
	s.log.Debug("resolving personal drive", logger.String("identity", identityID))

	d, err := s.vfs.GetPersonalDrive(ctx, "/personal")
	if err != nil {
		s.log.Error("failed to get personal drive", logger.Error(err))
		return Drive{}, fmt.Errorf("failed to get personal drive: %w", err)
	}

	return Drive{
		ID:       d.ID,
		Name:     d.Name,
		Type:     d.Type,
		Owner:    d.Owner,
		ReadOnly: d.ReadOnly,
	}, nil
}
