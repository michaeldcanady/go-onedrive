package drive

import (
	"context"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
)

// DriveSource defines the filesystem capabilities required by the drive service.
type DriveSource interface {
	ListDrives(ctx context.Context, provider string) ([]fs.Drive, error)
	GetPersonalDrive(ctx context.Context, provider string) (fs.Drive, error)
}

// MountProvider defines the interface for resolving mount information.
type MountProvider interface {
	ListMounts(ctx context.Context) ([]mount.MountConfig, error)
}

// Logger defines the interface required for logging within the drive service.
type Logger interface {
	Debug(msg string, fields ...logger.Field)
	Error(msg string, fields ...logger.Field)
}

// NewDefaultService initializes a new instance of the DefaultService.
func NewDefaultService(source DriveSource, mounts MountProvider, l Logger) *DefaultService {
	return &DefaultService{
		source: source,
		mounts: mounts,
		log:    l,
	}
}

// DefaultService provides the default implementation of the drive service.
type DefaultService struct {
	source DriveSource
	mounts MountProvider
	log    Logger
}

// ListDrives retrieves all accessible drives across all applicable mount points.
func (s *DefaultService) ListDrives(ctx context.Context, identityID string) ([]Drive, error) {
	s.log.Debug("listing drives", logger.String("identity", identityID))

	mountConfigs, err := s.mounts.ListMounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list mounts: %w", err)
	}

	var allDrives []Drive
	for _, m := range mountConfigs {
		if identityID != "" && m.IdentityID != identityID {
			continue
		}

		drives, err := s.source.ListDrives(ctx, m.Path)
		if err != nil {
			// Skip mounts that don't support drive discovery
			continue
		}

		for _, d := range drives {
			allDrives = append(allDrives, Drive{
				ID:       d.ID,
				Name:     d.Name,
				Type:     d.Type,
				Owner:    d.Owner,
				ReadOnly: d.ReadOnly,
			})
		}
	}

	return allDrives, nil
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

	drives, err := s.ListDrives(ctx, identityID)
	if err != nil {
		return Drive{}, err
	}

	for _, d := range drives {
		if d.Type == "personal" {
			return d, nil
		}
	}

	return Drive{}, fmt.Errorf("no personal drive found")
}
