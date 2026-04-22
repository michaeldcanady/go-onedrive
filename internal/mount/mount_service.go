package mount

import (
	"context"
	"fmt"
)

// ConfigRepository defines the interface for accessing mount configuration.
type ConfigRepository interface {
	GetMounts(ctx context.Context) ([]MountConfig, error)
	SaveMounts(ctx context.Context, mounts []MountConfig) error
}

// MountService is an implementation of the Service interface.
type MountService struct {
	configRepo ConfigRepository
}

// NewMountService creates a new instance of MountService.
func NewMountService(configRepo ConfigRepository) *MountService {
	return &MountService{
		configRepo: configRepo,
	}
}

// ListMounts retrieves all configured mount points.
func (s *MountService) ListMounts(ctx context.Context) ([]MountConfig, error) {
	return s.configRepo.GetMounts(ctx)
}

// AddMount adds or updates a mount point in the configuration.
func (s *MountService) AddMount(ctx context.Context, m MountConfig) error {
	mounts, err := s.configRepo.GetMounts(ctx)
	if err != nil {
		return err
	}

	found := false
	for i, existing := range mounts {
		if existing.Path == m.Path {
			mounts[i] = m
			found = true
			break
		}
	}

	if !found {
		mounts = append(mounts, m)
	}

	return s.configRepo.SaveMounts(ctx, mounts)
}

// RemoveMount removes a mount point from the configuration.
func (s *MountService) RemoveMount(ctx context.Context, path string) error {
	mounts, err := s.configRepo.GetMounts(ctx)
	if err != nil {
		return err
	}

	newMounts := make([]MountConfig, 0, len(mounts))
	for _, m := range mounts {
		if m.Path != path {
			newMounts = append(newMounts, m)
		}
	}

	if len(newMounts) == len(mounts) {
		return fmt.Errorf("mount point %s not found", path)
	}

	return s.configRepo.SaveMounts(ctx, newMounts)
}
