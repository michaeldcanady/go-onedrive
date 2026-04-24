package profile

import (
	"context"
	"fmt"
	"sync"
	"time"

	environment "github.com/michaeldcanady/go-onedrive/internal/core/env"
	"github.com/michaeldcanady/go-onedrive/internal/core/shared"
)

// DefaultService is an implementation of the profile.Service.
type DefaultService struct {
	profileRepo  ProfileRepository
	settingsRepo SettingsRepository
	repo         *BoltRepository
	env          environment.Service

	mu             sync.RWMutex
	sessionProfile string
}

// NewDefaultService initializes a new instance of the DefaultService.
func NewDefaultService(env environment.Service) (*DefaultService, error) {
	configDir, err := env.ConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	repo, err := NewBoltRepository(fmt.Sprintf("%s/profiles.db", configDir))
	if err != nil {
		return nil, err
	}

	s := &DefaultService{
		profileRepo:  repo,
		settingsRepo: repo,
		repo:         repo,
		env:          env,
	}

	return s, nil
}

// Bootstrap ensures the profile service is correctly initialized with default data.
func Bootstrap(ctx context.Context, s Service) error {
	exists, err := s.Exists(ctx, shared.DefaultProfileName)
	if err != nil {
		return err
	}

	if !exists {
		if _, err := s.Create(ctx, shared.DefaultProfileName); err != nil {
			return err
		}
		if err := s.SetActive(ctx, shared.DefaultProfileName, shared.ScopeGlobal); err != nil {
			return err
		}
	}

	return nil
}

// Close closes the database connection.
func (s *DefaultService) Close() error {
	return s.repo.Close()
}

// Get returns the profile with the specified name.
func (s *DefaultService) Get(ctx context.Context, name string) (Profile, error) {
	return s.profileRepo.Get(ctx, name)
}

// List returns a list of all profiles.
func (s *DefaultService) List(ctx context.Context) ([]Profile, error) {
	return s.profileRepo.List(ctx)
}

// Create generates a new profile with the specified name.
func (s *DefaultService) Create(ctx context.Context, name string) (Profile, error) {
	configDir, _ := s.env.ConfigDir()
	p := Profile{
		Name:       name,
		ConfigPath: fmt.Sprintf("%s/%s.yaml", configDir, name),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.profileRepo.Create(ctx, p); err != nil {
		return Profile{}, err
	}
	return p, nil
}

// Delete removes the profile with the specified name.
func (s *DefaultService) Delete(ctx context.Context, name string) error {
	return s.profileRepo.Delete(ctx, name)
}

// Exists checks if a profile with the specified name exists.
func (s *DefaultService) Exists(ctx context.Context, name string) (bool, error) {
	return s.profileRepo.Exists(ctx, name)
}

// Update saves the specified profile.
func (s *DefaultService) Update(ctx context.Context, p Profile) error {
	p.UpdatedAt = time.Now()
	return s.profileRepo.Update(ctx, p)
}

// GetActive retrieves the currently active profile.
func (s *DefaultService) GetActive(ctx context.Context) (Profile, error) {
	s.mu.RLock()
	name := s.sessionProfile
	s.mu.RUnlock()

	var err error
	// Fallback to settings repository if no session override
	if name == "" {
		name, err = s.settingsRepo.GetSetting(ctx, "active_profile")
		if err != nil {
			// Default to shared.DefaultProfileName if not set
			name = shared.DefaultProfileName
		}
	}

	return s.Get(ctx, name)
}

// SetActive marks a specific profile as the active one.
func (s *DefaultService) SetActive(ctx context.Context, name string, scope shared.Scope) error {
	exists, err := s.profileRepo.Exists(ctx, name)
	if err != nil {
		return err
	}
	if !exists {
		return ErrProfileNotFound
	}

	if scope == shared.ScopeSession {
		s.mu.Lock()
		s.sessionProfile = name
		s.mu.Unlock()
		return nil
	}

	return s.settingsRepo.SetSetting(ctx, "active_profile", name)
}
