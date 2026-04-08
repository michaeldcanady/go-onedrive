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
func (s *DefaultService) ListDrives(ctx context.Context) ([]Drive, error) {
	log := s.log.WithContext(ctx)
	log.Debug("listing drives")

	drives, err := s.gateway.ListDrives(ctx)
	if err != nil {
		log.Error("failed to list drives", logger.Error(err))
		return nil, fmt.Errorf("failed to list drives: %w", err)
	}

	return drives, nil
}

// ResolveDrive identifies a drive by its ID or name.
func (s *DefaultService) ResolveDrive(ctx context.Context, driveRef string) (Drive, error) {
	log := s.log.WithContext(ctx).With(logger.String("ref", driveRef))
	log.Debug("resolving drive")

	drives, err := s.ListDrives(ctx)
	if err != nil {
		log.Error("failed to resolve drive due to list failure", logger.Error(err))
		return Drive{}, err
	}

	for _, d := range drives {
		if strings.EqualFold(d.ID, driveRef) || strings.EqualFold(d.Name, driveRef) {
			log.Debug("drive resolved successfully", logger.String("id", d.ID), logger.String("name", d.Name))
			return d, nil
		}
	}

	log.Warn("drive not found during resolution")
	return Drive{}, fmt.Errorf("drive %s not found", driveRef)
}

// ResolvePersonalDrive retrieves the user's primary personal drive.
func (s *DefaultService) ResolvePersonalDrive(ctx context.Context) (Drive, error) {
	log := s.log.WithContext(ctx)
	log.Debug("resolving personal drive")

	d, err := s.gateway.GetPersonalDrive(ctx)
	if err != nil {
		log.Error("failed to get personal drive", logger.Error(err))
		return Drive{}, fmt.Errorf("failed to get personal drive: %w", err)
	}

	log.Debug("personal drive resolved successfully", logger.String("id", d.ID))
	return d, nil
}

// GetActive retrieves the currently active drive.
func (s *DefaultService) GetActive(ctx context.Context) (Drive, error) {
	log := s.log.WithContext(ctx)
	id, err := s.state.Get(state.KeyDrive)
	if err != nil {
		log.Error("failed to get active drive ID from state", logger.Error(err))
		return Drive{}, fmt.Errorf("failed to get active drive ID: %w", err)
	}

	if id == "" {
		log.Debug("no active drive set, falling back to personal drive")
		return s.ResolvePersonalDrive(ctx)
	}

	log.Debug("retrieving active drive", logger.String("id", id))
	return s.ResolveDrive(ctx, id)
}

// SetActive marks a specific drive as the active one with the given scope.
func (s *DefaultService) SetActive(ctx context.Context, driveID string, scope state.Scope) error {
	log := s.log.WithContext(ctx).With(logger.String("id", driveID), logger.String("scope", scope.String()))
	log.Info("setting active drive")
	return s.state.Set(state.KeyDrive, driveID, scope)
}
