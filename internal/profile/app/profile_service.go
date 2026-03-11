package app

import (
	"context"
	"errors"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	logger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	domainprofile "github.com/michaeldcanady/go-onedrive/internal/profile/domain"
)

// Ensure interface compliance
var _ domainprofile.ProfileService = (*ProfileService)(nil)

type ProfileService struct {
	log  logger.Logger
	repo domainprofile.Repository
}

func New(l logger.Logger, repo domainprofile.Repository) *ProfileService {
	return &ProfileService{
		log:  l,
		repo: repo,
	}
}

const (
	eventProfileCreateStart   = "profiledomain.create.start"
	eventProfileCreateSuccess = "profiledomain.create.success"
	eventProfileCreateFailure = "profiledomain.create.failure"

	eventProfileDeleteStart   = "profiledomain.delete.start"
	eventProfileDeleteSuccess = "profiledomain.delete.success"
	eventProfileDeleteFailure = "profiledomain.delete.failure"

	eventProfileExistsStart   = "profiledomain.exists.start"
	eventProfileExistsSuccess = "profiledomain.exists.success"
	eventProfileExistsFailure = "profiledomain.exists.failure"

	eventProfileListStart   = "profiledomain.list.start"
	eventProfileListSuccess = "profiledomain.list.success"
	eventProfileListFailure = "profiledomain.list.failure"

	eventProfileGetRepoLoad = "profiledomain.get.repo.load"
	eventProfileGetSuccess  = "profiledomain.get.success"
	eventProfileGetFailure  = "profiledomain.get.failure"
)

func (s *ProfileService) Create(ctx context.Context, name string) (domainprofile.Profile, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domainprofile.Profile{}, errors.New("profile name is empty")
	}

	correlationID := util.CorrelationIDFromContext(ctx)

	log := s.log.WithContext(ctx).With(
		logger.String("correlation_id", correlationID),
		logger.String("profile_name", name),
		logger.String("event", eventProfileCreateStart),
	)

	log.Info("creating profile")

	profile, err := s.repo.Create(name)
	if err != nil {
		log.Error("failed to create profile",
			logger.String("event", eventProfileCreateFailure),
			logger.Error(err),
		)
		return domainprofile.Profile{}, err
	}

	log.Info("profile created successfully",
		logger.String("event", eventProfileCreateSuccess),
	)

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

	log := s.log.WithContext(ctx).With(
		logger.String("correlation_id", correlationID),
		logger.String("profile_name", name),
		logger.String("event", eventProfileDeleteStart),
	)

	log.Info("deleting profile")

	if err := s.repo.Delete(name); err != nil {
		log.Error("failed to delete profile",
			logger.String("event", eventProfileDeleteFailure),
			logger.Error(err),
		)
		return err
	}

	log.Info("profile deleted successfully",
		logger.String("event", eventProfileDeleteSuccess),
	)

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

	log := s.log.WithContext(ctx).With(
		logger.String("correlation_id", correlationID),
		logger.String("profile_name", name),
		logger.String("event", eventProfileExistsStart),
	)

	log.Debug("checking if profile exists")

	exists, err := s.repo.Exists(name)
	if err != nil {
		log.Error("failed to check profile existence",
			logger.String("event", eventProfileExistsFailure),
			logger.Error(err),
		)
		return false, err
	}

	log.Debug("profile existence check complete",
		logger.String("event", eventProfileExistsSuccess),
		logger.Bool("exists", exists),
	)

	return exists, nil
}

// ───────────────────────────────────────────────────────────────────────────────
// List
// ───────────────────────────────────────────────────────────────────────────────

func (s *ProfileService) List(ctx context.Context) ([]domainprofile.Profile, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	log := s.log.WithContext(ctx).With(
		logger.String("correlation_id", correlationID),
		logger.String("event", eventProfileListStart),
	)

	log.Info("listing profiles")

	list, err := s.repo.List()
	if err != nil {
		log.Error("failed to list profiles",
			logger.String("event", eventProfileListFailure),
			logger.Error(err),
		)
		return nil, err
	}

	log.Info("profile list retrieved",
		logger.String("event", eventProfileListSuccess),
		logger.Int("count", len(list)),
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

	log := s.log.WithContext(ctx).With(
		logger.String("correlation_id", correlationID),
		logger.String("profile_name", name),
	)

	name = strings.TrimSpace(name)
	if name == "" {
		log.Error("profile name is empty",
			logger.String("event", eventProfileGetFailure),
		)
		return profile, errors.New("profile name is empty")
	}

	log.Info("retrieving profile")

	// Load from repository
	profile, err := s.repo.Get(ctx, name)
	if err != nil {
		log.Error("failed to load profile from repository",
			logger.String("event", eventProfileGetFailure),
			logger.Error(err),
		)
		return profile, err
	}

	log.Info("profile loaded from repository",
		logger.String("event", eventProfileGetRepoLoad),
	)

	log.Info("profile retrieval complete",
		logger.String("event", eventProfileGetSuccess),
	)

	return profile, nil
}
