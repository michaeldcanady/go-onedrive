package cliprofileservicego

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	domainprofile "github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/core"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
)

type ProfileService struct {
	cacheSvc   CacheService
	logger     logging.Logger
	profileDir string
}

func New(cacheSvc CacheService, logger logging.Logger, profileDir string) *ProfileService {
	return &ProfileService{
		cacheSvc:   cacheSvc,
		logger:     logger,
		profileDir: profileDir,
	}
}

func (s *ProfileService) GetProfile(ctx context.Context, name string) (domainprofile.Profile, error) {
	var profile domainprofile.Profile

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

func (s *ProfileService) resolveProfile(ctx context.Context, name string) (domainprofile.Profile, error) {
	var profile domainprofile.Profile

	if strings.TrimSpace(s.profileDir) == "" {
		return profile, errors.New("missing profile directory")
	}

	name = strings.ToLower(name)

	profile.Name = name
	profile.Directory = filepath.Join(s.profileDir, name)
	profile.ConfigurationPath = filepath.Join(profile.Directory)

	return profile, nil
}
