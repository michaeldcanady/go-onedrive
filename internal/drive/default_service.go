package drive

import (
	"context"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/state"
)

// DefaultService provides the default implementation of the drive service.
type DefaultService struct {
	gateway Gateway
	state   state.Service
	log     logger.Logger
}

// NewDefaultService initializes a new instance of the DefaultService.
func NewDefaultService(gateway Gateway, state state.Service, l logger.Logger) *DefaultService {
	return &DefaultService{
		gateway: gateway,
		state:   state,
		log:     l,
	}
}

// ListDrives retrieves all accessible OneDrive drives.
func (s *DefaultService) ListDrives(ctx context.Context, identityID string) ([]Drive, error) {
	s.log.Debug("listing drives", logger.String("identity", identityID))

	drives, err := s.gateway.ListDrives(ctx, identityID)
	if err != nil {
		s.log.Error("failed to list drives", logger.Error(err))
		return nil, fmt.Errorf("failed to list drives: %w", err)
	}

	return drives, nil
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

	d, err := s.gateway.GetPersonalDrive(ctx, identityID)
	if err != nil {
		s.log.Error("failed to get personal drive", logger.Error(err))
		return Drive{}, fmt.Errorf("failed to get personal drive: %w", err)
	}

	return d, nil
}

// GetActive retrieves the currently active drive.
func (s *DefaultService) GetActive(ctx context.Context, identityID string) (Drive, error) {
	var id string
	var err error

	if identityID != "" {
		id, err = s.state.GetScoped("tokens/microsoft", identityID+"/active_drive")
	} else {
		id, err = s.state.Get(state.KeyDrive)
	}

	if err != nil {
		if err != state.ErrKeyNotFound {
			return Drive{}, fmt.Errorf("failed to get active drive ID: %w", err)
		}
		// If not found, fall back to personal drive
		id = ""
	}

	if id == "" {
		return s.ResolvePersonalDrive(ctx, identityID)
	}

	return s.ResolveDrive(ctx, id, identityID)
}

// SetActive marks a specific drive as the active one with the given scope.
func (s *DefaultService) SetActive(ctx context.Context, driveID string, identityID string, scope state.Scope) error {
	if identityID != "" {
		return s.state.SetScoped("tokens/microsoft", identityID+"/active_drive", driveID, scope)
	}
	return s.state.Set(state.KeyDrive, driveID, scope)
}
