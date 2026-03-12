package drive

import (
	"context"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
)

// DefaultService provides the default implementation of the drive service.
type DefaultService struct {
	gateway Gateway
	log     logger.Logger
}

// NewDefaultService initializes a new instance of the DefaultService.
func NewDefaultService(gateway Gateway, l logger.Logger) *DefaultService {
	return &DefaultService{
		gateway: gateway,
		log:     l,
	}
}

// ListDrives retrieves all accessible OneDrive drives.
func (s *DefaultService) ListDrives(ctx context.Context) ([]Drive, error) {
	s.log.Debug("listing drives")

	drives, err := s.gateway.ListDrives(ctx)
	if err != nil {
		s.log.Error("failed to list drives", logger.Error(err))
		return nil, fmt.Errorf("failed to list drives: %w", err)
	}

	return drives, nil
}

// ResolveDrive identifies a drive by its ID or name.
func (s *DefaultService) ResolveDrive(ctx context.Context, driveRef string) (Drive, error) {
	s.log.Debug("resolving drive", logger.String("ref", driveRef))

	drives, err := s.ListDrives(ctx)
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
func (s *DefaultService) ResolvePersonalDrive(ctx context.Context) (Drive, error) {
	s.log.Debug("resolving personal drive")

	d, err := s.gateway.GetPersonalDrive(ctx)
	if err != nil {
		s.log.Error("failed to get personal drive", logger.Error(err))
		return Drive{}, fmt.Errorf("failed to get personal drive: %w", err)
	}

	return d, nil
}
