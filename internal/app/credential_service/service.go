package credentialservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/event"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
)

const (
	maxAuthAttempts = 3
)

type Service struct {
	profileService       ProfileService
	configurationService ConfigurationService
	credential           azcore.TokenCredential
	logger               logging.Logger
	publisher            event.Publisher
}

func New(profileSvc ProfileService, configurationService ConfigurationService, publisher event.Publisher, logger logging.Logger) *Service {
	s := &Service{
		profileService:       profileSvc,
		configurationService: configurationService,
		publisher:            publisher,
		logger:               logger,
	}

	return s
}

func (s *Service) LoadCredential(ctx context.Context, profile *azidentity.AuthenticationRecord) (azcore.TokenCredential, error) {
	// If already loaded, return cached credential
	if s.credential != nil {
		s.logger.Debug("credential already loaded; returning cached instance")
		return s.credential, nil
	}

	authType, err := s.configurationService.GetString(ctx, "auth.type")
	if err != nil {
		s.logger.Error("unable to get auth.type", logging.Any("error", err))
		return nil, errors.Join(errors.New("unable to get auth.type"), err)
	}
	clientID, err := s.configurationService.GetString(ctx, "auth.client_id")
	if err != nil {
		s.logger.Error("unable to get auth.client_id", logging.Any("error", err))
		return nil, errors.Join(errors.New("unable to get auth.client_id"), err)
	}
	tenantID, err := s.configurationService.GetString(ctx, "auth.tenant_id")
	if err != nil {
		s.logger.Error("unable to get auth.tenant_id", logging.Any("error", err))
		return nil, errors.Join(errors.New("unable to get auth.tenant_id"), err)
	}
	redirectURI, err := s.configurationService.GetString(ctx, "auth.redirect_uri")
	if err != nil {
		s.logger.Error("unable to get auth.redirect_uri", logging.Any("error", err))
		return nil, errors.Join(errors.New("unable to get auth.redirect_uri"), err)
	}

	authCfg := config.AuthenticationConfigImpl{
		Type:                 authType,
		ClientID:             clientID,
		TenantID:             tenantID,
		RedirectURI:          redirectURI,
		AuthenticationRecord: profile,
	}

	s.logger.Debug("authentication config loaded", logging.Any("config", authCfg))

	if profile == nil {
		s.logger.Info("no cached authentication profile found")
	} else {
		s.logger.Info("loaded cached authentication profile")
		s.logger.Debug("cached authentication profile", logging.Any("profile", profile))
	}

	authCfg.AuthenticationRecord = profile

	// Create credential
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

func (s *Service) Authenticate(ctx context.Context, opts ...AuthenticationOption) (azcore.AccessToken, error) {
	config := &authenticationConfig{}
	buildConfig(config, opts...)

	var (
		token   azcore.AccessToken
		success bool = false
	)

	for range maxAuthAttempts {
		profile, err := s.profileService.Load(ctx)
		if err != nil {
			return azcore.AccessToken{}, fmt.Errorf("unable to load profile: %w", err)
		}

		cred, err := s.LoadCredential(ctx, profile)
		if err != nil {
			return azcore.AccessToken{}, fmt.Errorf("unable to load credential: %w", err)
		}

		if (profile == nil || isEmptyRecord(*profile)) || config.Force {
			s.logger.Info("Starting authentication flow...")

			authenticator, ok := cred.(Authenticator)
			if !ok {
				return azcore.AccessToken{}, errors.New("configured credential does not support explicit authentication")
			}

			record, err := authenticator.Authenticate(ctx, &policy.TokenRequestOptions{
				Scopes:    config.Scopes,
				EnableCAE: config.EnableCAE,
				Claims:    config.Claims,
			})
			if err != nil {
				s.logger.Error("authentication failed", logging.String("error", err.Error()))
				return azcore.AccessToken{}, fmt.Errorf("authentication failed: %w", err)
			}
			profile := &record

			if err := s.profileService.Save(ctx, profile); err != nil {
				return azcore.AccessToken{}, fmt.Errorf("unable to save profile: %w", err)
			}
		}

		token, err = cred.GetToken(ctx, policy.TokenRequestOptions{
			Scopes:    config.Scopes,
			EnableCAE: config.EnableCAE,
			Claims:    config.Claims,
		})
		if err != nil {
			if isAuthRequired(err) {
				s.logger.Info("authentication required; retrying...")
				profile = nil
				continue
			}
			return azcore.AccessToken{}, fmt.Errorf("unable to acquire token: %w", err)
		}
		success = true
		break
	}

	if !success {
		return azcore.AccessToken{}, errors.New("maximum authentication attempts exceeded")
	}

	return token, nil
}

func (s *Service) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	profile, err := s.profileService.Load(ctx)
	if err != nil {
		return azcore.AccessToken{}, fmt.Errorf("unable to load profile: %w", err)
	}

	cred, err := s.LoadCredential(ctx, profile)
	if err != nil {
		return azcore.AccessToken{}, fmt.Errorf("unable to load credential: %w", err)
	}

	return cred.GetToken(ctx, options)
}

func isAuthRequired(err error) bool {
	var authErr *azidentity.AuthenticationRequiredError
	return errors.As(err, &authErr)
}
