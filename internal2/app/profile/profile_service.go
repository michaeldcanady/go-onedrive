package cliprofileservicego

import (
	"context"
	"errors"
	"strings"

	domainprofile "github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/core"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
)

var _ domainprofile.ProfileService = (*ProfileService)(nil)

type ProfileService struct {
	cacheSvc CacheService
	logger   logging.Logger
	repo     domainprofile.ProfileService
}

func New(cacheSvc CacheService, logger logging.Logger, repo domainprofile.ProfileService) *ProfileService {
	return &ProfileService{
		cacheSvc: cacheSvc,
		logger:   logger,
		repo:     repo,
	}
}

func (s *ProfileService) Create(name string) (domainprofile.Profile, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domainprofile.Profile{}, errors.New("profile name is empty")
	}

	profile, err := s.repo.Create(name)
	if err != nil {
		return domainprofile.Profile{}, err
	}

	// Cache it
	if err := s.cacheSvc.SetCLIProfile(context.Background(), name, profile); err != nil {
		s.logger.Warn("failed to cache profile", logging.String("err", err.Error()))
	}

	return profile, nil
}

func (s *ProfileService) Delete(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("profile name is empty")
	}

	if err := s.repo.Delete(name); err != nil {
		return err
	}

	// Remove from cache
	//if err := s.cacheSvc.DeleteCLIProfile(context.Background(), name); err != nil {
	//	s.logger.Warn("failed to delete cached profile", logging.String("err", err.Error()))
	//}

	return nil
}

func (s *ProfileService) Exists(name string) (bool, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return false, errors.New("profile name is empty")
	}
	return s.repo.Exists(name)
}

func (s *ProfileService) List() ([]domainprofile.Profile, error) {
	return s.repo.List()
}

func (s *ProfileService) Get(ctx context.Context, name string) (domainprofile.Profile, error) {
	return s.GetProfile(ctx, name)
}

func (s *ProfileService) GetProfile(ctx context.Context, name string) (domainprofile.Profile, error) {
	var profile domainprofile.Profile

	if err := ctx.Err(); err != nil {
		return profile, err
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return profile, errors.New("profile name is empty")
	}

	// Try cache first
	profile, err := s.cacheSvc.GetCLIProfile(ctx, name)
	if err == nil && profile != (domainprofile.Profile{}) {
		return profile, nil
	}

	if err != core.ErrKeyNotFound {
		return profile, err
	}

	// Cache miss → load from repository
	profile, err = s.repo.Get(ctx, name)
	if err != nil {
		return profile, err
	}

	// Cache it
	if err := s.cacheSvc.SetCLIProfile(ctx, name, profile); err != nil {
		return profile, errors.Join(errors.New("unable to cache profile"), err)
	}

	return profile, nil
}
