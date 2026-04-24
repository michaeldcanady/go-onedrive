package mount

import (
	"context"
	"fmt"
	"sync"
)

// MountService is an implementation of the Service interface.
type MountService struct {
	configRepo          ConfigRepository
	mu                  sync.RWMutex
	validators          map[string]OptionValidator
	completionProviders map[string]CompletionProvider
}

// NewMountService creates a new instance of MountService.
func NewMountService(configRepo ConfigRepository) *MountService {
	return &MountService{
		configRepo:          configRepo,
		validators:          make(map[string]OptionValidator),
		completionProviders: make(map[string]CompletionProvider),
	}
}

// RegisterValidator registers a validator for a given mount type.
func (s *MountService) RegisterValidator(mountType string, v OptionValidator) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.validators[mountType] = v
}

// RegisterCompletionProvider registers a completion provider for a given mount type.
func (s *MountService) RegisterCompletionProvider(mountType string, p CompletionProvider) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.completionProviders[mountType] = p
}

// GetCompletionProvider retrieves a registered completion provider.
func (s *MountService) GetCompletionProvider(mountType string) (CompletionProvider, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.completionProviders[mountType]
	return p, ok
}

// GetMountOptions retrieves all registered mount options.
func (s *MountService) GetMountOptions() map[string][]MountOption {
	s.mu.RLock()
	defer s.mu.RUnlock()
	options := make(map[string][]MountOption)
	for mountType, validator := range s.validators {
		if provider, ok := validator.(OptionsProvider); ok {
			options[mountType] = provider.ProvideOptions()
		}
	}
	return options
}

// ListMounts retrieves all configured mount points.
func (s *MountService) ListMounts(ctx context.Context) ([]MountConfig, error) {
	return s.configRepo.GetMounts(ctx)
}

// AddMount adds or updates a mount point in the configuration.
func (s *MountService) AddMount(ctx context.Context, m MountConfig) error {
	s.mu.RLock()
	validator, ok := s.validators[m.Type]
	s.mu.RUnlock()

	if ok {
		if err := validator.ValidateOptions(m.Options); err != nil {
			return fmt.Errorf("invalid options for mount type %s: %w", m.Type, err)
		}
	}

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
