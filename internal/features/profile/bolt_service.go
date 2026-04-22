package profile

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/core/env"
	"github.com/michaeldcanady/go-onedrive/internal/core/shared"
	bolt "go.etcd.io/bbolt"
)

// DefaultService is an implementation of the profile.Service.
type DefaultService struct {
	profileRepo  ProfileRepository
	settingsRepo SettingsRepository
	env          environment.Service
	db           *bolt.DB // Kept for cleanup

	mu             sync.RWMutex
	sessionProfile string
}

// NewDefaultService initializes a new instance of the DefaultService.
func NewDefaultService(env environment.Service) (*DefaultService, error) {
	configDir, err := env.ConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	dbFilePath := filepath.Join(configDir, "profiles.db")
	db, err := bolt.Open(dbFilePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open BoltDB: %w", err)
	}

	// Ensure buckets are created
	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("profiles")); err != nil {
			return fmt.Errorf("failed to create profiles bucket: %w", err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte("settings")); err != nil {
			return fmt.Errorf("failed to create settings bucket: %w", err)
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}

	repo := NewBoltRepository(db)
	s := &DefaultService{
		profileRepo:  repo,
		settingsRepo: repo,
		env:          env,
		db:           db,
	}

	// Ensure default profile exists
	exists, _ := repo.Exists(context.Background(), shared.DefaultProfileName)
	if !exists {
		_, _ = s.Create(context.Background(), shared.DefaultProfileName)
		_ = repo.SetSetting(context.Background(), "active_profile", shared.DefaultProfileName)
	}

	return s, nil
}

// ResolvePath returns the configuration file path for the specified profile name.
func (s *DefaultService) ResolvePath(ctx context.Context, profileName string) (string, error) {
	p, err := s.profileRepo.Get(ctx, profileName)
	if err != nil {
		return "", err
	}
	return p.ConfigPath, nil
}

// Close closes the BoltDB database connection.
func (s *DefaultService) Close() error {
	return s.db.Close()
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
		ConfigPath: filepath.Join(configDir, fmt.Sprintf("%s.yaml", name)),
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
