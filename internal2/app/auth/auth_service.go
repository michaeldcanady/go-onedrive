package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"

	authdomain "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	authinfra "github.com/michaeldcanady/go-onedrive/internal2/infra/auth"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
)

var _ authdomain.AuthService = (*AuthService)(nil)

type AuthService struct {
	cache   CacheService
	factory authinfra.CredentialFactory
	config  ConfigurationService
	logger  logging.Logger
}

func NewService(
	factory authinfra.CredentialFactory,
	cache CacheService,
	config ConfigurationService,
	logger logging.Logger,
) *AuthService {
	return &AuthService{
		factory: factory,
		cache:   cache,
		config:  config,
		logger:  logger,
	}
}

// GetToken implements auth.AuthService.
// This method NEVER triggers interactive login.
// It only attempts silent token acquisition using cached AuthenticationRecord.
func (s *AuthService) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	s.logger.Debug("attempting silent token acquisition")

	// Load cached record
	record, err := s.cache.GetProfile(ctx, "default")
	if err != nil {
		return azcore.AccessToken{}, fmt.Errorf("failed to load cached authentication record: %w", err)
	}

	// Load profile configuration
	cfg, err := s.config.GetConfiguration(ctx, "default")
	if err != nil {
		return azcore.AccessToken{}, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Build credential options
	credOpts := authinfra.CredentialOptions{
		Type:                 cfg.Auth.Type,
		ClientID:             cfg.Auth.ClientID,
		TenantID:             cfg.Auth.TenantID,
		AuthenticationRecord: record,
	}

	// Create credential
	cred, err := s.factory.Create(credOpts)
	if err != nil {
		return azcore.AccessToken{}, fmt.Errorf("failed to create credential: %w", err)
	}

	// Attempt silent token acquisition
	token, err := cred.GetToken(ctx, options)
	if err != nil {
		if record == (azidentity.AuthenticationRecord{}) {
			return azcore.AccessToken{}, errors.New("not logged in")
		}
		return azcore.AccessToken{}, fmt.Errorf("failed to acquire token silently: %w", err)
	}

	s.logger.Debug("silent token acquisition successful")
	return token, nil
}

// Login performs interactive authentication if needed, then retrieves a token.
func (s *AuthService) Login(ctx context.Context, profileName string, opts authdomain.LoginOptions) (*authdomain.LoginResult, error) {
	s.logger.Info("starting login flow", logging.String("profile", profileName))

	var (
		record azidentity.AuthenticationRecord
		err    error
	)

	// Load cached record (may be empty)
	if s.cache != nil {
		record, err = s.cache.GetProfile(ctx, profileName)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve cached auth record: %w", err)
		}
	}

	// Load profile configuration
	cfg, err := s.config.GetConfiguration(ctx, profileName)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve configuration: %w", err)
	}

	// Build credential options
	credOpts := authinfra.CredentialOptions{
		Type:                 cfg.Auth.Type,
		ClientID:             cfg.Auth.ClientID,
		TenantID:             cfg.Auth.TenantID,
		AuthenticationRecord: record,
	}

	// Create credential
	cred, err := s.factory.Create(credOpts)
	if err != nil {
		s.logger.Error("failed to retrieve credential", logging.String("error", err.Error()))
		return nil, fmt.Errorf("failed to retrieve credential: %w", err)
	}

	authenticator, ok := cred.(authdomain.Authenticator)
	if !ok {
		return nil, errors.New("credential does not support explicit authentication")
	}

	tokenOpts := &policy.TokenRequestOptions{
		Scopes:    opts.Scopes,
		EnableCAE: opts.EnableCAE,
	}

	var token azcore.AccessToken
	maxAttempts := 3

	for range maxAttempts {
		// Determine if interactive login is required
		needsAuth := opts.Force || record == (azidentity.AuthenticationRecord{})

		if needsAuth {
			s.logger.Info("performing interactive authentication")

			newRecord, err := authenticator.Authenticate(ctx, tokenOpts)
			if err != nil {
				return nil, fmt.Errorf("authentication failed: %w", err)
			}

			record = newRecord

			// Reload credential with updated record
			cred, err = s.factory.Create(authinfra.CredentialOptions{
				Type:                 cfg.Auth.Type,
				ClientID:             cfg.Auth.ClientID,
				TenantID:             cfg.Auth.TenantID,
				AuthenticationRecord: newRecord,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to reload credential: %w", err)
			}
		}

		// Retrieve token
		token, err = cred.GetToken(ctx, *tokenOpts)
		if err != nil {
			if !isAuthRequired(err) {
				return nil, fmt.Errorf("failed to retrieve token: %w", err)
			}
			record = azidentity.AuthenticationRecord{}
			continue
		}
		break
	}

	if token == (azcore.AccessToken{}) {
		s.logger.Warn("token is empty")
		return nil, errors.New("empty token")
	}

	// Save updated record
	if err := s.cache.SetProfile(ctx, profileName, record); err != nil {
		s.logger.Warn("failed to save authentication record", logging.String("error", err.Error()))
	}

	s.logger.Info("login successful")

	return &authdomain.LoginResult{
		AccessToken: token.Token,
		RecordSaved: (record != azidentity.AuthenticationRecord{}),
	}, nil
}
