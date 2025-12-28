package app

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/spf13/viper"
)

type CredentialService interface {
	LoadCredential(ctx context.Context) (azcore.TokenCredential, error)
}

type CredentialServiceImpl struct {
	profileService ProfileService
}

func NewCredentialService(profileSvc ProfileService) *CredentialServiceImpl {
	return &CredentialServiceImpl{profileService: profileSvc}
}

func (s *CredentialServiceImpl) LoadCredential(ctx context.Context) (azcore.TokenCredential, error) {
	sub := viper.Sub("auth")
	if sub == nil {
		return nil, fmt.Errorf("missing 'auth' config section")
	}

	var authCfg config.AuthenticationConfigImpl
	if err := sub.Unmarshal(&authCfg); err != nil {
		return nil, fmt.Errorf("unable to unmarshal auth config: %w", err)
	}

	// Load cached profile
	record, err := s.profileService.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load cached profile: %w", err)
	}

	authCfg.AuthenticationRecord = record

	cred, err := CredentialFactory(&authCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	return cred, nil
}
