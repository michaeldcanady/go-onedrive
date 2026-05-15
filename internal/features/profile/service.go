package profile

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
)

type profileService struct {
	repo   Repository
	logger logger.Service
}

// NewProfileService returns a new [Service] initialized with the provided repository.
func NewProfileService(repo Repository, l logger.Service) Service {
	return &profileService{
		repo:   repo,
		logger: l,
	}
}

func (s *profileService) Create(name string) (*Profile, error) {
	p := &Profile{Name: name}
	if err := s.repo.Create(p); err != nil {
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}
	s.logger.Info("profile created", "name", name)
	return p, nil
}

func (s *profileService) List() ([]*Profile, error) {
	return s.repo.List()
}

func (s *profileService) Delete(name string) error {
	if err := s.repo.Delete(name); err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}
	s.logger.Info("profile deleted", "name", name)
	return nil
}

func (s *profileService) GetCurrent() (*Profile, error) {
	name, err := s.repo.GetCurrent()
	if err != nil {
		return nil, fmt.Errorf("failed to get current profile name: %w", err)
	}
	if name == "" {
		return nil, nil
	}
	return &Profile{Name: name}, nil
}

func (s *profileService) SetCurrent(name string) error {
	// Optional: verify profile exists before setting as current
	if err := s.repo.SetCurrent(name); err != nil {
		return fmt.Errorf("failed to set current profile: %w", err)
	}
	s.logger.Info("current profile set", "name", name)
	return nil
}
