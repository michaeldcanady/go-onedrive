package cliprofileservicego

import (
	"context"
	"errors"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	domainprofile "github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/core"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

var _ domainprofile.ProfileService = (*ProfileService)(nil)

type ProfileService struct {
	cacheSvc cache.CacheService
	logger   logging.Logger
	repo     domainprofile.ProfileService
}

func New(cacheSvc cache.CacheService, logger logging.Logger, repo domainprofile.ProfileService) *ProfileService {
	return &ProfileService{
		cacheSvc: cacheSvc,
		logger:   logger,
		repo:     repo,
	}
}

func (s *ProfileService) Create(name string) (domainprofile.Profile, error) {
	cid := util.CorrelationIDFromContext(context.Background())

	name = strings.TrimSpace(name)
	if name == "" {
		s.logger.Warn("profile name is empty",
			logging.String("event", "cli_profile_create_invalid"),
			logging.String("correlation_id", cid),
		)
		return domainprofile.Profile{}, errors.New("profile name is empty")
	}

	s.logger.Debug("creating profile",
		logging.String("event", "cli_profile_create_start"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	profile, err := s.repo.Create(name)
	if err != nil {
		s.logger.Error("failed to create profile in repository",
			logging.String("event", "cli_profile_create_repo_error"),
			logging.String("profile", name),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return domainprofile.Profile{}, err
	}

	// Cache it
	if err := s.cacheSvc.SetCLIProfile(context.Background(), name, profile); err != nil {
		s.logger.Warn("failed to cache profile",
			logging.String("event", "cli_profile_create_cache_error"),
			logging.String("profile", name),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
	}

	s.logger.Info("profile created successfully",
		logging.String("event", "cli_profile_create_success"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	return profile, nil
}

func (s *ProfileService) Delete(name string) error {
	cid := util.CorrelationIDFromContext(context.Background())

	name = strings.TrimSpace(name)
	if name == "" {
		s.logger.Warn("profile name is empty",
			logging.String("event", "cli_profile_delete_invalid"),
			logging.String("correlation_id", cid),
		)
		return errors.New("profile name is empty")
	}

	s.logger.Debug("deleting profile",
		logging.String("event", "cli_profile_delete_start"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	if err := s.repo.Delete(name); err != nil {
		s.logger.Error("failed to delete profile from repository",
			logging.String("event", "cli_profile_delete_repo_error"),
			logging.String("profile", name),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return err
	}

	// Optional: remove from cache
	// (Your cache service does not currently expose DeleteCLIProfile)
	// If added later, logging should follow the same conventions.

	s.logger.Info("profile deleted successfully",
		logging.String("event", "cli_profile_delete_success"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	return nil
}

func (s *ProfileService) Exists(name string) (bool, error) {
	cid := util.CorrelationIDFromContext(context.Background())

	name = strings.TrimSpace(name)
	if name == "" {
		s.logger.Warn("profile name is empty",
			logging.String("event", "cli_profile_exists_invalid"),
			logging.String("correlation_id", cid),
		)
		return false, errors.New("profile name is empty")
	}

	s.logger.Debug("checking if profile exists",
		logging.String("event", "cli_profile_exists_start"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	exists, err := s.repo.Exists(name)
	if err != nil {
		s.logger.Error("failed to check profile existence",
			logging.String("event", "cli_profile_exists_repo_error"),
			logging.String("profile", name),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return false, err
	}

	s.logger.Info("profile existence check completed",
		logging.String("event", "cli_profile_exists_success"),
		logging.String("profile", name),
		logging.Bool("exists", exists),
		logging.String("correlation_id", cid),
	)

	return exists, nil
}

func (s *ProfileService) List() ([]domainprofile.Profile, error) {
	cid := util.CorrelationIDFromContext(context.Background())

	s.logger.Debug("listing profiles",
		logging.String("event", "cli_profile_list_start"),
		logging.String("correlation_id", cid),
	)

	list, err := s.repo.List()
	if err != nil {
		s.logger.Error("failed to list profiles",
			logging.String("event", "cli_profile_list_repo_error"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return nil, err
	}

	s.logger.Info("profiles listed successfully",
		logging.String("event", "cli_profile_list_success"),
		logging.Int("count", len(list)),
		logging.String("correlation_id", cid),
	)

	return list, nil
}

func (s *ProfileService) Get(ctx context.Context, name string) (domainprofile.Profile, error) {
	return s.GetProfile(ctx, name)
}

func (s *ProfileService) GetProfile(ctx context.Context, name string) (domainprofile.Profile, error) {
	cid := util.CorrelationIDFromContext(ctx)

	var profile domainprofile.Profile

	if err := ctx.Err(); err != nil {
		s.logger.Warn("context canceled while retrieving profile",
			logging.String("event", "cli_profile_get_ctx_error"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return profile, err
	}

	name = strings.TrimSpace(name)
	if name == "" {
		s.logger.Warn("profile name is empty",
			logging.String("event", "cli_profile_get_invalid"),
			logging.String("correlation_id", cid),
		)
		return profile, errors.New("profile name is empty")
	}

	s.logger.Debug("retrieving profile",
		logging.String("event", "cli_profile_get_start"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	// Try cache first
	profile, err := s.cacheSvc.GetCLIProfile(ctx, name)
	if err == nil && profile != (domainprofile.Profile{}) {
		s.logger.Info("profile retrieved from cache",
			logging.String("event", "cli_profile_get_cache_hit"),
			logging.String("profile", name),
			logging.String("correlation_id", cid),
		)
		return profile, nil
	}

	if err != nil && err != core.ErrKeyNotFound {
		s.logger.Error("failed to retrieve profile from cache",
			logging.String("event", "cli_profile_get_cache_error"),
			logging.String("profile", name),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return profile, err
	}

	s.logger.Debug("cache miss; loading profile from repository",
		logging.String("event", "cli_profile_get_cache_miss"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	profile, err = s.repo.Get(ctx, name)
	if err != nil {
		s.logger.Error("failed to retrieve profile from repository",
			logging.String("event", "cli_profile_get_repo_error"),
			logging.String("profile", name),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return profile, err
	}

	// Cache it
	if err := s.cacheSvc.SetCLIProfile(ctx, name, profile); err != nil {
		s.logger.Warn("failed to cache profile",
			logging.String("event", "cli_profile_get_cache_set_error"),
			logging.String("profile", name),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
	}

	s.logger.Info("profile retrieved successfully",
		logging.String("event", "cli_profile_get_success"),
		logging.String("profile", name),
		logging.String("correlation_id", cid),
	)

	return profile, nil
}
