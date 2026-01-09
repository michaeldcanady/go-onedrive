package cacheservice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal/event"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
)

const (
	homeAccountIDKey = "homeAccountID"
	errorKey         = "error"
)

// Service provides methods to manage cache entries.
type Service struct {
	publisher event.Publisher
	logger    logging.Logger
	cache     Cache
}

// New creates a new instance of the cache service.
func New(cache Cache, publisher event.Publisher, logger logging.Logger) *Service {
	return &Service{
		publisher: publisher,
		logger:    logger,
		cache:     cache,
	}
}

// PutProfile stores an authentication profile in the cache and returns its HomeAccountID.
func (s *Service) PutProfile(ctx context.Context, profile *azidentity.AuthenticationRecord) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	id := profile.HomeAccountID

	data, err := json.Marshal(profile)
	if err != nil {
		return "", fmt.Errorf("failed to marshal profile: %w", err)
	}

	old, err := s.GetProfile(ctx, id)
	if err != nil {
		s.logger.Error("failed to get old profile",
			logging.String(homeAccountIDKey, id),
			logging.Any(errorKey, err),
		)
		return "", fmt.Errorf("failed to get old profile: %w", err)
	}

	if err := s.cache.Put(ctx, id, data); err != nil {
		s.logger.Error("failed to put profile in cache",
			logging.String(homeAccountIDKey, id),
			logging.Any(errorKey, err),
		)
		return "", fmt.Errorf("failed to cache profile: %w", err)
	}

	if s.publisher == nil {
		s.logger.Warn("no event publisher configured; skipping profile.updated event")
	} else {
		s.logger.Debug("publishing profile.updated event")
		if err := s.publisher.Publish(ctx, newProfileUpdatedEvent(old, profile)); err != nil {
			s.logger.Error("failed to publish profile.updated event",
				logging.String(homeAccountIDKey, id),
				logging.Any(errorKey, err),
			)
			return "", fmt.Errorf("failed to publish profile.updated event: %w", err)
		}
	}

	return id, nil
}

// GetProfile retrieves an authentication profile from the cache by its HomeAccountID.
func (s *Service) GetProfile(ctx context.Context, id string) (*azidentity.AuthenticationRecord, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	entry, err := s.cache.Get(ctx, id)
	if err != nil {
		s.logger.Error("failed to get profile from cache",
			logging.String(homeAccountIDKey, id),
			logging.Any(errorKey, err),
		)
		return nil, fmt.Errorf("failed to get profile from cache: %w", err)
	}
	if entry == nil {
		return nil, nil
	}

	var profile azidentity.AuthenticationRecord
	if err := json.Unmarshal(entry, &profile); err != nil {
		s.logger.Error("failed to unmarshal profile from cache",
			logging.String(homeAccountIDKey, id),
			logging.Any(errorKey, err),
		)
		return nil, fmt.Errorf("failed to unmarshal profile: %w", err)
	}

	return &profile, nil
}

// DeleteProfile removes an authentication profile from the cache by its HomeAccountID.
func (s *Service) DeleteProfile(ctx context.Context, id string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	old, err := s.GetProfile(ctx, id)
	if err != nil {
		s.logger.Error("failed to get old profile",
			logging.String(homeAccountIDKey, id),
			logging.Any(errorKey, err),
		)
		return fmt.Errorf("failed to get old profile: %w", err)
	}

	if err := s.cache.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete profile from cache",
			logging.String(homeAccountIDKey, id),
			logging.Any(errorKey, err),
		)
		return fmt.Errorf("failed to delete profile from cache: %w", err)
	}

	if s.publisher == nil {
		s.logger.Warn("no event publisher configured; skipping profile.deleted event")
	} else {
		s.logger.Debug("publishing profile.deleted event")
		if err := s.publisher.Publish(ctx, newProfileDeletedEvent(old)); err != nil {
			s.logger.Error("failed to publish profile.deleted event",
				logging.String(homeAccountIDKey, id),
				logging.Any(errorKey, err),
			)
			return fmt.Errorf("failed to publish profile.deleted event: %w", err)
		}
	}

	return nil
}
