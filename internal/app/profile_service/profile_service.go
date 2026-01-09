package profileservice

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal/cache/abstractions"
	"github.com/michaeldcanady/go-onedrive/internal/event"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
)

const profileKey = "profile.json"

type Service struct {
	store     abstractions.Store
	codec     abstractions.Codec
	publisher event.Publisher
	logger    logging.Logger

	cached *azidentity.AuthenticationRecord
}

func New(store abstractions.Store, codec abstractions.Codec, publisher event.Publisher, logger logging.Logger) *Service {
	return &Service{
		store:     store,
		codec:     codec,
		publisher: publisher,
		logger:    logger,
	}
}

func (s *Service) Load(ctx context.Context) (*azidentity.AuthenticationRecord, error) {
	// Return cached profile if available
	if s.cached != nil {
		s.logger.Debug("returning cached profile")
		return s.cached, nil
	}

	data, err := s.store.LoadBytes(ctx, profileKey)
	if err != nil {
		s.logger.Error("failed to load profile bytes", logging.Any("error", err))
		return nil, fmt.Errorf("load profile bytes: %w", err)
	}

	if data == nil {
		s.logger.Debug("profile not found in store")
		return nil, nil
	}

	var p azidentity.AuthenticationRecord
	if err := s.codec.Decode(data, &p); err != nil {
		s.logger.Error("failed to decode profile", logging.Any("error", err))
		return nil, fmt.Errorf("decode profile: %w", err)
	}

	if isZeroProfile(&p) {
		s.logger.Debug("profile exists but is zero/invalid")
		return nil, nil
	}

	s.cached = &p
	s.logger.Info("profile loaded successfully")

	// ❌ DO NOT publish profile.loaded here anymore
	// Load() is not a state change

	return &p, nil
}

func (s *Service) Save(ctx context.Context, p *azidentity.AuthenticationRecord) error {
	if p == nil || isZeroProfile(p) {
		s.logger.Debug("saving zero profile; writing empty record")
		p = &azidentity.AuthenticationRecord{}
	}

	data, err := s.codec.Encode(p)
	if err != nil {
		s.logger.Error("failed to encode profile", logging.Any("error", err))
		return fmt.Errorf("encode profile: %w", err)
	}

	if err := s.store.SaveBytes(ctx, profileKey, data); err != nil {
		s.logger.Error("failed to save profile bytes", logging.Any("error", err))
		return fmt.Errorf("save profile bytes: %w", err)
	}

	s.cached = p

	// Publish profile.saved
	if s.publisher != nil {
		s.logger.Debug("publishing profile.saved event")
		if err := s.publisher.Publish(ctx, newProfileSavedEvent(*p)); err != nil {
			s.logger.Warn("failed to publish profile.saved event", logging.Any("error", err))
		}
	}

	s.logger.Info("profile saved successfully")
	return nil
}

func (s *Service) Clear(ctx context.Context) error {
	// Save an empty profile
	err := s.Save(ctx, nil)
	if err != nil {
		s.logger.Error("failed to clear profile", logging.Any("error", err))
		return err
	}

	s.cached = nil

	s.logger.Info("profile cleared")

	// Publish profile.cleared
	if s.publisher != nil {
		s.logger.Debug("publishing profile.cleared event")
		if err := s.publisher.Publish(ctx, newProfileClearedEvent()); err != nil {
			s.logger.Warn("failed to publish profile.cleared event", logging.Any("error", err))
			return fmt.Errorf("unable to publish profile.cleared event: %w", err)
		}
	}

	return nil
}

func isZeroProfile(p *azidentity.AuthenticationRecord) bool {
	return p == nil
}
