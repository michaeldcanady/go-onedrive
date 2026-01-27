package auth

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal/event"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/core"
)

type Service struct {
	cacheService  CacheService
	credential    azcore.TokenCredential
	logger        logging.Logger
	publisher     event.Publisher
	configService ConfigurationService
}

func New(cacheService CacheService, publisher event.Publisher, logger logging.Logger, configService ConfigurationService) *Service {
	s := &Service{
		cacheService:  cacheService,
		publisher:     publisher,
		logger:        logger,
		configService: configService,
	}

	return s
}

func (s *Service) LoadCredential(ctx context.Context, profileName string) (azcore.TokenCredential, error) {
	if s.credential != nil {
		s.logger.Debug("credential already loaded; returning cached instance")
		return s.credential, nil
	}

	// 1. Load profile configuration
	cfg, err := s.configService.GetConfiguration(ctx, profileName)
	if err != nil {
		s.logger.Error("unable to load profile configuration", logging.Any("error", err))
		return nil, err
	}

	// 2. Extract authentication config
	authCfg := cfg.Auth

	// 3. Load cached authentication record
	var record azidentity.AuthenticationRecord
	if s.cacheService != nil {
		record, err = s.cacheService.GetProfile(ctx, profileName)
		if err != nil && !errors.Is(err, core.ErrKeyNotFound) {
			s.logger.Error("unable to load cached authentication record", logging.Any("error", err))
			return nil, err
		}
	} else {
		s.logger.Warn("missing caching service")
	}

	if record == (azidentity.AuthenticationRecord{}) {
		s.logger.Info("no cached authentication record found")
		authCfg.AuthenticationRecord = nil
	} else {
		s.logger.Info("loaded cached authentication record")
		authCfg.AuthenticationRecord = &record
	}

	// 4. Create credential
	s.credential, err = CredentialFactory(&authCfg)
	if err != nil {
		s.logger.Error("unable to create credential", logging.Any("error", err))
		return nil, err
	}

	if s.cacheService != nil {
		if err := s.cacheService.SetProfile(ctx, profileName, *authCfg.AuthenticationRecord); err != nil {
			s.logger.Warn("unable to cache profile", logging.Any("error", err))
		} else {
			s.logger.Info("cached profile")
		}
	} else {
		s.logger.Warn("authentication profile not cached")
	}

	s.logger.Info("credential created successfully")

	// 5. Publish event
	if s.publisher != nil {
		s.logger.Debug("publishing credential.loaded event")
		if err := s.publisher.Publish(ctx, newCredentialLoadedEvent(s.credential)); err != nil {
			s.logger.Error("failed to publish credential.loaded event", logging.Any("error", err))
		}
	}

	return s.credential, nil
}
