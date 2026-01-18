package credentialservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/event"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	"github.com/spf13/viper"
)

const (
	defaultProfileKey = "default"
)

type Service struct {
	cacheService CacheService
	credential   azcore.TokenCredential
	logger       logging.Logger
	publisher    event.Publisher
}

func New(cacheService CacheService, publisher event.Publisher, logger logging.Logger) *Service {
	s := &Service{
		cacheService: cacheService,
		publisher:    publisher,
		logger:       logger,
	}

	return s
}

func (s *Service) LoadCredential(ctx context.Context) (azcore.TokenCredential, error) {
	// If already loaded, return cached credential
	if s.credential != nil {
		s.logger.Debug("credential already loaded; returning cached instance")
		return s.credential, nil
	}

	sub := viper.Sub("auth")
	if sub == nil {
		err := fmt.Errorf("missing 'auth' config section")
		s.logger.Error(err.Error())
		return nil, err
	}

	var authCfg config.AuthenticationConfigImpl
	if err := sub.Unmarshal(&authCfg); err != nil {
		s.logger.Error("unable to unmarshal auth config", logging.Any("error", err))
		return nil, errors.Join(errors.New("unable to unmarshal auth config"), err)
	}
	s.logger.Debug("authentication config loaded", logging.Any("config", authCfg))

	// Load cached profile
	record, err := s.cacheService.GetProfile(ctx, defaultProfileKey)
	if err != nil {
		s.logger.Error("unable to load cached authentication record", logging.Any("error", err))
		return nil, errors.Join(errors.New("unable to load cached authentication record"), err)
	}

	if record == (azidentity.AuthenticationRecord{}) {
		s.logger.Info("no cached authentication record found")
		authCfg.AuthenticationRecord = nil
	} else {
		s.logger.Info("loaded cached authentication record")
		s.logger.Debug("cached authentication record", logging.Any("record", record))
		authCfg.AuthenticationRecord = &record
	}

	// TODO: never gets updated profile back
	s.credential, err = CredentialFactory(&authCfg)
	if err != nil {
		s.logger.Error("unable to create credential", logging.Any("error", err))
		return nil, errors.Join(errors.New("unable to create credential"), err)
	}

	s.logger.Info("credential created successfully")
	s.logger.Debug("credential instance", logging.Any("credential", s.credential))

	// Publish event
	if s.publisher != nil {
		s.logger.Debug("publishing credential.loaded event")
		if err := s.publisher.Publish(ctx, newCredentialLoadedEvent(s.credential)); err != nil {
			s.logger.Error("failed to publish credential.loaded event", logging.Any("error", err))
		}
	}

	return s.credential, nil
}
