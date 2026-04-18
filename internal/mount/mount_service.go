package mount

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/config"
)

// MountService is an implementation of the Service interface.
type MountService struct {
	configSvc config.Service
}

// NewMountService creates a new instance of MountService.
func NewMountService(configSvc config.Service) *MountService {
	return &MountService{
		configSvc: configSvc,
	}
}

// ListMounts retrieves all configured mount points.
func (s *MountService) ListMounts(ctx context.Context) ([]config.MountConfig, error) {
	cfg, err := s.configSvc.GetConfig(ctx)
	if err != nil {
		return nil, err
	}
	return cfg.Mounts, nil
}

// AddMount adds or updates a mount point in the configuration.
func (s *MountService) AddMount(ctx context.Context, m config.MountConfig) error {
	cfg, err := s.configSvc.GetConfig(ctx)
	if err != nil {
		return err
	}

	found := false
	for i, existing := range cfg.Mounts {
		if existing.Path == m.Path {
			cfg.Mounts[i] = m
			found = true
			break
		}
	}

	if !found {
		cfg.Mounts = append(cfg.Mounts, m)
	}

	return s.configSvc.SaveConfig(ctx, cfg)
}

// RemoveMount removes a mount point from the configuration.
func (s *MountService) RemoveMount(ctx context.Context, path string) error {
	cfg, err := s.configSvc.GetConfig(ctx)
	if err != nil {
		return err
	}

	newMounts := make([]config.MountConfig, 0, len(cfg.Mounts))
	for _, m := range cfg.Mounts {
		if m.Path != path {
			newMounts = append(newMounts, m)
		}
	}

	if len(newMounts) == len(cfg.Mounts) {
		return fmt.Errorf("mount point %s not found", path)
	}

	cfg.Mounts = newMounts
	return s.configSvc.SaveConfig(ctx, cfg)
}
