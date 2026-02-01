// internal2/infra/profile/fs_profile_service.go
package profile

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
)

const (
	DefaultProfileName = "default"
)

type FSProfileService struct {
	baseDir string // e.g. ~/.config/odc
}

func NewFSProfileService(baseDir string) *FSProfileService {
	return &FSProfileService{baseDir: baseDir}
}

func (s *FSProfileService) profilePath(name string) string {
	return filepath.Join(s.baseDir, name)
}

func (s *FSProfileService) Exists(name string) (bool, error) {
	if err := s.ensureDefault(); err != nil {
		return false, err
	}

	_, err := os.Stat(s.profilePath(name))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *FSProfileService) List() ([]profile.Profile, error) {
	if err := s.ensureDefault(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, err
	}

	var profiles []profile.Profile
	for _, e := range entries {
		if e.IsDir() {
			profiles = append(profiles, profile.Profile{
				Name: e.Name(),
				Path: s.profilePath(e.Name()),
			})
		}
	}
	return profiles, nil
}

func (s *FSProfileService) Create(name string) (profile.Profile, error) {
	if err := s.ensureDefault(); err != nil {
		return profile.Profile{}, err
	}

	p := s.profilePath(name)
	if err := os.MkdirAll(p, 0o700); err != nil {
		return profile.Profile{}, err
	}
	return profile.Profile{Name: name, Path: p}, nil
}

func (s *FSProfileService) Delete(name string) error {
	if name == DefaultProfileName {
		return errors.New("default profile can't be deleted")
	}
	p := s.profilePath(name)
	return os.RemoveAll(p)
}

func (s *FSProfileService) Get(_ context.Context, name string) (profile.Profile, error) {
	if err := s.ensureDefault(); err != nil {
		return profile.Profile{}, err
	}

	exists, err := s.Exists(name)
	if err != nil {
		return profile.Profile{}, err
	}
	if !exists {
		return profile.Profile{}, fmt.Errorf("profile %q does not exist", name)
	}
	return profile.Profile{Name: name, Path: s.profilePath(name)}, nil
}

func (s *FSProfileService) ensureDefault() error {
	defaultProfilePath := filepath.Join(s.baseDir, DefaultProfileName)
	return os.MkdirAll(defaultProfilePath, os.ModePerm)
}
