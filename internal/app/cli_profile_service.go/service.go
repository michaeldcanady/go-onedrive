package cliprofileservicego

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/cachev2/core"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
)

type Service struct {
	cacheSvc   CacheService
	logger     logging.Logger
	profileDir string
}

func New(cacheSvc CacheService, logger logging.Logger, profileDir string) *Service {
	return &Service{
		cacheSvc:   cacheSvc,
		logger:     logger,
		profileDir: profileDir,
	}
}

func (s *Service) GetProfile(ctx context.Context, name string) (Profile, error) {
	var profile Profile

	if err := ctx.Err(); err != nil {
		return profile, err
	}

	if strings.TrimSpace(name) == "" {
		return profile, errors.New("name is empty")
	}

	profile, err := s.cacheSvc.GetCLIProfile(ctx, name)
	if err == nil {
		return profile, nil
	}

	if err != core.ErrKeyNotFound {
		return profile, err
	}

	profile, err = s.resolveProfile(ctx, name)
	if err != nil {
		return profile, err
	}

	if err := s.cacheSvc.SetCLIProfile(ctx, name, profile); err != nil {
		return profile, errors.Join(errors.New("unable to cache profile"), err)
	}

	return profile, nil
}

func (s *Service) resolveProfile(ctx context.Context, name string) (Profile, error) {
	var profile Profile

	if strings.TrimSpace(s.profileDir) == "" {
		return profile, errors.New("missing profile directory")
	}

	name = strings.ToLower(name)

	profile.Name = name
	profile.Directory = filepath.Join(s.profileDir, name)
	profile.ConfigurationPath = filepath.Join(profile.Directory)

	return profile, nil
}
