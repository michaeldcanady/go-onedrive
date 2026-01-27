package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal/cachev2/core"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
)

type Cache interface {
	GetProfile(ctx context.Context, profile string) (azidentity.AuthenticationRecord, error)
	SetProfile(ctx context.Context, profile string, record azidentity.AuthenticationRecord) error
}

type CredentialProvider struct {
	cache   Cache
	config  ConfigProvider
	factory CredentialFactory
	log     logging.Logger
}

func NewCredentialProvider(
	cache Cache,
	config ConfigProvider,
	factory CredentialFactory,
	log logging.Logger,
) *CredentialProvider {
	return &CredentialProvider{
		cache:   cache,
		config:  config,
		factory: factory,
		log:     log,
	}
}

func (p *CredentialProvider) Credential(ctx context.Context, profile string) (azcore.TokenCredential, error) {
	cfg, err := p.config.GetConfiguration(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf("load profile config: %w", err)
	}

	authCfg := cfg.Auth

	// Load cached auth record
	var record azidentity.AuthenticationRecord
	if p.cache != nil {
		record, err = p.cache.GetProfile(ctx, profile)
		if err != nil && !errors.Is(err, core.ErrKeyNotFound) {
			return nil, fmt.Errorf("load cached auth record: %w", err)
		}
		if record != (azidentity.AuthenticationRecord{}) {
			authCfg.AuthenticationRecord = &record
		}
	}

	// Create credential via factory
	cred, err := p.factory.Create(&authCfg)
	if err != nil {
		return nil, fmt.Errorf("create credential: %w", err)
	}

	// Save updated auth record
	if p.cache != nil && authCfg.AuthenticationRecord != nil {
		if err := p.cache.SetProfile(ctx, profile, *authCfg.AuthenticationRecord); err != nil {
			p.log.Warn("failed to cache auth record", logging.Any("error", err))
		}
	}

	return cred, nil
}
