package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	"github.com/spf13/viper"
)

type CredentialService interface {
	LoadCredential(ctx context.Context) (azcore.TokenCredential, error)
}

type CredentialServiceImpl struct {
	profileService ProfileService
	credential     azcore.TokenCredential
	logger         logging.Logger
}

func NewCredentialService(profileSvc ProfileService, logger logging.Logger) *CredentialServiceImpl {
	return &CredentialServiceImpl{
		profileService: profileSvc,
		logger:         logger,
	}
}

func (s *CredentialServiceImpl) LoadCredential(ctx context.Context) (azcore.TokenCredential, error) {
	sub := viper.Sub("auth")
	if sub == nil {
		s.logger.Error("missing 'auth' config section")
		return nil, fmt.Errorf("missing 'auth' config section")
	}

	var authCfg config.AuthenticationConfigImpl
	if err := sub.Unmarshal(&authCfg); err != nil {
		s.logger.Error("unable to unmarshal auth config")
		return nil, errors.Join(errors.New("unable to unmarshal auth config"), err)
	}

	// Load cached profile
	record, err := s.profileService.Load(ctx)
	if err != nil {
		s.logger.Error("unable to load cached authentication record")
		return nil, errors.Join(errors.New("unable to load cached authentication record"), err)
	}
	s.logger.Info("loaded cached authentication record")
	s.logger.Debug("cached authentication record", logging.Any("record", record))

	authCfg.AuthenticationRecord = record

	s.credential, err = CredentialFactory(&authCfg)
	if err != nil {
		s.logger.Error("unable to create credential")
		return nil, errors.Join(errors.New("unable to create credential"), err)
	}
	s.logger.Info("created credential")
	s.logger.Debug("credential", logging.Any("credential", s.credential))

	return s.credential, nil
}
