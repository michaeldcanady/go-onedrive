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

// Ensure interface compliance
var _ domainprofile.ProfileService = (*ProfileService)(nil)

type ProfileService struct {
	cacheSvc cache.CacheService
	logger   logging.Logger
	repo     domainprofile.Repository
}

func New(cacheSvc cache.CacheService, logger logging.Logger, repo domainprofile.Repository) *ProfileService {
	return &ProfileService{
		cacheSvc: cacheSvc,
		logger:   logger,
		repo:     repo,
	}
}

const (
	eventProfileCreateStart   = "profile.create.start"
	eventProfileCreateSuccess = "profile.create.success"
	eventProfileCreateFailure = "profile.create.failure"
	eventProfileCreateCache   = "profile.create.cache"

	eventProfileDeleteStart   = "profile.delete.start"
	eventProfileDeleteSuccess = "profile.delete.success"
	eventProfileDeleteFailure = "profile.delete.failure"

	eventProfileExistsStart   = "profile.exists.start"
	eventProfileExistsSuccess = "profile.exists.success"
	eventProfileExistsFailure = "profile.exists.failure"

	eventProfileListStart   = "profile.list.start"
	eventProfileListSuccess = "profile.list.success"
	eventProfileListFailure = "profile.list.failure"

	eventProfileGetStart      = "profile.get.start"
	eventProfileGetCacheHit   = "profile.get.cache.hit"
	eventProfileGetCacheMiss  = "profile.get.cache.miss"
	eventProfileGetRepoLoad   = "profile.get.repo.load"
	eventProfileGetCacheSave  = "profile.get.cache.save"
	eventProfileGetCacheError = "profile.get.cache.error"
	eventProfileGetSuccess    = "profile.get.success"
	eventProfileGetFailure    = "profile.get.failure"
)

func (s *ProfileService) Create(ctx context.Context, name string) (domainprofile.Profile, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domainprofile.Profile{}, errors.New("profile name is empty")
	}

	correlationID := util.CorrelationIDFromContext(ctx)

	logger := s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("profile_name", name),
		logging.String("event", eventProfileCreateStart),
	)

	logger.Info("creating profile")

	profile, err := s.repo.Create(name)
	if err != nil {
		logger.Error("failed to create profile",
			logging.String("event", eventProfileCreateFailure),
			logging.Error(err),
		)
		return domainprofile.Profile{}, err
	}

	logger.Info("profile created successfully",
		logging.String("event", eventProfileCreateSuccess),
	)

	// Cache it
	if err := s.cacheSvc.SetCLIProfile(ctx, name, profile); err != nil {
		logger.Warn("failed to cache profile",
			logging.String("event", eventProfileCreateCache),
			logging.Error(err),
		)
	}

	return profile, nil
}

// ───────────────────────────────────────────────────────────────────────────────
// Delete
// ───────────────────────────────────────────────────────────────────────────────

func (s *ProfileService) Delete(ctx context.Context, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("profile name is empty")
	}

	correlationID := util.CorrelationIDFromContext(ctx)

	logger := s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("profile_name", name),
		logging.String("event", eventProfileDeleteStart),
	)

	logger.Info("deleting profile")

	if err := s.repo.Delete(name); err != nil {
		logger.Error("failed to delete profile",
			logging.String("event", eventProfileDeleteFailure),
			logging.Error(err),
		)
		return err
	}

	logger.Info("profile deleted successfully",
		logging.String("event", eventProfileDeleteSuccess),
	)

	// Optional: delete from cache
	// if err := s.cacheSvc.DeleteCLIProfile(ctx, name); err != nil {
	//     logger.Warn("failed to delete cached profile", logging.Error(err))
	// }

	return nil
}

// ───────────────────────────────────────────────────────────────────────────────
// Exists
// ───────────────────────────────────────────────────────────────────────────────

func (s *ProfileService) Exists(ctx context.Context, name string) (bool, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return false, errors.New("profile name is empty")
	}

	correlationID := util.CorrelationIDFromContext(ctx)

	logger := s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("profile_name", name),
		logging.String("event", eventProfileExistsStart),
	)

	logger.Debug("checking if profile exists")

	exists, err := s.repo.Exists(name)
	if err != nil {
		logger.Error("failed to check profile existence",
			logging.String("event", eventProfileExistsFailure),
			logging.Error(err),
		)
		return false, err
	}

	logger.Debug("profile existence check complete",
		logging.String("event", eventProfileExistsSuccess),
		logging.Bool("exists", exists),
	)

	return exists, nil
}

// ───────────────────────────────────────────────────────────────────────────────
// List
// ───────────────────────────────────────────────────────────────────────────────

func (s *ProfileService) List(ctx context.Context) ([]domainprofile.Profile, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	logger := s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("event", eventProfileListStart),
	)

	logger.Info("listing profiles")

	list, err := s.repo.List()
	if err != nil {
		logger.Error("failed to list profiles",
			logging.String("event", eventProfileListFailure),
			logging.Error(err),
		)
		return nil, err
	}

	logger.Info("profile list retrieved",
		logging.String("event", eventProfileListSuccess),
		logging.Int("count", len(list)),
	)

	return list, nil
}

// ───────────────────────────────────────────────────────────────────────────────
// Get / GetProfile
// ───────────────────────────────────────────────────────────────────────────────

func (s *ProfileService) Get(ctx context.Context, name string) (domainprofile.Profile, error) {
	return s.GetProfile(ctx, name)
}

func (s *ProfileService) GetProfile(ctx context.Context, name string) (domainprofile.Profile, error) {
	var profile domainprofile.Profile

	if err := ctx.Err(); err != nil {
		return profile, err
	}

	correlationID := util.CorrelationIDFromContext(ctx)

	logger := s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("profile_name", name),
	)

	name = strings.TrimSpace(name)
	if name == "" {
		logger.Error("profile name is empty",
			logging.String("event", eventProfileGetFailure),
		)
		return profile, errors.New("profile name is empty")
	}

	logger.Info("retrieving profile")

	// Try cache first
	profile, err := s.cacheSvc.GetCLIProfile(ctx, name)
	if err == nil && profile != (domainprofile.Profile{}) {
		logger.Info("profile retrieved from cache",
			logging.String("event", eventProfileGetCacheHit),
		)
		return profile, nil
	}

	if err != nil && err != core.ErrKeyNotFound {
		logger.Error("failed to retrieve profile from cache",
			logging.String("event", eventProfileGetFailure),
			logging.Error(err),
		)
		return profile, err
	}

	logger.Info("profile not found in cache",
		logging.String("event", eventProfileGetCacheMiss),
	)

	// Load from repository
	profile, err = s.repo.Get(ctx, name)
	if err != nil {
		logger.Error("failed to load profile from repository",
			logging.String("event", eventProfileGetFailure),
			logging.Error(err),
		)
		return profile, err
	}

	logger.Info("profile loaded from repository",
		logging.String("event", eventProfileGetRepoLoad),
	)

	// Cache it
	if err := s.cacheSvc.SetCLIProfile(ctx, name, profile); err != nil {
		logger.Warn("failed to cache profile",
			logging.String("event", eventProfileGetCacheError),
			logging.Error(err),
		)
		return profile, errors.Join(errors.New("unable to cache profile"), err)
	}

	logger.Info("profile cached successfully",
		logging.String("event", eventProfileGetCacheSave),
	)

	logger.Info("profile retrieval complete",
		logging.String("event", eventProfileGetSuccess),
	)

	return profile, nil
}
