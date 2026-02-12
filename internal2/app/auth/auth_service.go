package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"

	authdomain "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/config"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/state"
	authinfra "github.com/michaeldcanady/go-onedrive/internal2/infra/auth"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

var _ authdomain.AuthService = (*AuthService)(nil)

const (
	eventAuthSilentStart       = "auth.silent.start"
	eventAuthSilentSuccess     = "auth.silent.success"
	eventAuthSilentFailure     = "auth.silent.failure"
	eventAuthSilentNotLoggedIn = "auth.silent.not_logged_in"

	eventAuthLoginStart         = "auth.login.start"
	eventAuthLoginRecordLoaded  = "auth.login.record_loaded"
	eventAuthLoginConfigLoaded  = "auth.login.config_loaded"
	eventAuthLoginInteractive   = "auth.login.interactive"
	eventAuthLoginTokenAttempt  = "auth.login.token_attempt"
	eventAuthLoginTokenSuccess  = "auth.login.token_success"
	eventAuthLoginTokenFailure  = "auth.login.token_failure"
	eventAuthLoginEmptyToken    = "auth.login.empty_token"
	eventAuthLoginRecordSaved   = "auth.login.record_saved"
	eventAuthLoginRecordSaveErr = "auth.login.record_save_error"
	eventAuthLoginSuccess       = "auth.login.success"

	eventAuthLogoutStart         = "auth.logout.start"
	eventAuthLogoutRecordLoad    = "auth.logout.record_load"
	eventAuthLogoutRecordMissing = "auth.logout.record_missing"
	eventAuthLogoutDelete        = "auth.logout.delete"
	eventAuthLogoutDeleteError   = "auth.logout.delete_error"
	eventAuthLogoutSuccess       = "auth.logout.success"
)

type AuthService struct {
	cache   cache.CacheService
	factory authinfra.CredentialFactory
	config  config.ConfigService
	state   state.Service
	logger  logging.Logger
}

func NewService(
	factory authinfra.CredentialFactory,
	cache cache.CacheService,
	config config.ConfigService,
	state state.Service,
	logger logging.Logger,
) *AuthService {
	return &AuthService{
		factory: factory,
		cache:   cache,
		config:  config,
		logger:  logger,
		state:   state,
	}
}

// GetToken implements auth.AuthService.
// This method NEVER triggers interactive login.
// It only attempts silent token acquisition using cached AuthenticationRecord.
func (s *AuthService) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	logger := s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("event", eventAuthSilentStart),
		logging.String("scopes", fmt.Sprintf("%v", options.Scopes)),
	)

	logger.Info("starting silent token acquisition")

	profileName, err := s.state.GetCurrentProfile()
	if err != nil {
		logger.Error("failed to determine current profile", logging.Error(err))
		return azcore.AccessToken{}, err
	}

	logger = logger.With(logging.String("profile", profileName))

	record, err := s.cache.GetProfile(ctx, profileName)
	if err != nil {
		logger.Error("failed to load cached authentication record", logging.Error(err))
		return azcore.AccessToken{}, fmt.Errorf("failed to load cached authentication record: %w", err)
	}

	logger.Debug("loaded cached authentication record",
		logging.Bool("has_record", record != (azidentity.AuthenticationRecord{})),
	)

	cfg, err := s.config.GetConfiguration(ctx, profileName)
	if err != nil {
		logger.Error("failed to load configuration", logging.Error(err))
		return azcore.AccessToken{}, fmt.Errorf("failed to load configuration: %w", err)
	}

	credOpts := authinfra.CredentialOptions{
		Type:                 cfg.Auth.Type,
		ClientID:             cfg.Auth.ClientID,
		TenantID:             cfg.Auth.TenantID,
		AuthenticationRecord: record,
	}

	logger.Debug("creating credential",
		logging.String("auth_type", credOpts.Type),
		logging.String("tenant_id", credOpts.TenantID),
		logging.String("client_id", credOpts.ClientID),
	)

	cred, err := s.factory.Create(credOpts)
	if err != nil {
		logger.Error("failed to create credential", logging.Error(err))
		return azcore.AccessToken{}, fmt.Errorf("failed to create credential: %w", err)
	}

	token, err := cred.GetToken(ctx, options)
	if err != nil {
		if record == (azidentity.AuthenticationRecord{}) {
			logger.Warn("silent token acquisition failed: no authentication record",
				logging.String("event", eventAuthSilentNotLoggedIn),
			)
			return azcore.AccessToken{}, errors.New("not logged in")
		}

		logger.Error("silent token acquisition failed",
			logging.String("event", eventAuthSilentFailure),
			logging.Error(err),
		)
		return azcore.AccessToken{}, fmt.Errorf("failed to acquire token silently: %w", err)
	}

	logger.Info("silent token acquisition successful",
		logging.String("event", eventAuthSilentSuccess),
		logging.String("expires_on", token.ExpiresOn.String()),
	)

	return token, nil
}

// Login performs interactive authentication if needed, then retrieves a token.
func (s *AuthService) Login(ctx context.Context, profileName string, opts authdomain.LoginOptions) (*authdomain.LoginResult, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	logger := s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("profile", profileName),
		logging.String("event", eventAuthLoginStart),
	)

	logger.Info("starting login flow",
		logging.Bool("force", opts.Force),
		logging.Bool("enable_cae", opts.EnableCAE),
		logging.String("scopes", fmt.Sprintf("%v", opts.Scopes)),
	)

	var (
		record azidentity.AuthenticationRecord
		err    error
	)

	record, err = s.cache.GetProfile(ctx, profileName)
	if err != nil {
		logger.Error("failed to retrieve cached auth record", logging.Error(err))
		return nil, fmt.Errorf("failed to retrieve cached auth record: %w", err)
	}

	logger.Debug("loaded cached authentication record",
		logging.String("event", eventAuthLoginRecordLoaded),
		logging.Bool("has_record", record != (azidentity.AuthenticationRecord{})),
	)

	cfg, err := s.config.GetConfiguration(ctx, profileName)
	if err != nil {
		logger.Error("failed to retrieve configuration", logging.Error(err))
		return nil, fmt.Errorf("failed to retrieve configuration: %w", err)
	}

	logger.Debug("configuration loaded",
		logging.String("event", eventAuthLoginConfigLoaded),
		logging.String("auth_type", cfg.Auth.Type),
		logging.String("tenant_id", cfg.Auth.TenantID),
		logging.String("client_id", cfg.Auth.ClientID),
	)

	credOpts := authinfra.CredentialOptions{
		Type:                 cfg.Auth.Type,
		ClientID:             cfg.Auth.ClientID,
		TenantID:             cfg.Auth.TenantID,
		AuthenticationRecord: record,
	}

	cred, err := s.factory.Create(credOpts)
	if err != nil {
		logger.Error("failed to create credential", logging.Error(err))
		return nil, fmt.Errorf("failed to retrieve credential: %w", err)
	}

	authenticator, ok := cred.(authdomain.Authenticator)
	if !ok {
		logger.Error("credential does not support explicit authentication")
		return nil, errors.New("credential does not support explicit authentication")
	}

	tokenOpts := &policy.TokenRequestOptions{
		Scopes:    opts.Scopes,
		EnableCAE: opts.EnableCAE,
	}

	var token azcore.AccessToken
	maxAttempts := 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		needsAuth := opts.Force || record == (azidentity.AuthenticationRecord{})

		if needsAuth {
			logger.Info("performing interactive authentication",
				logging.String("event", eventAuthLoginInteractive),
				logging.Int("attempt", attempt),
				logging.Int("max_attempts", maxAttempts),
			)

			newRecord, err := authenticator.Authenticate(ctx, tokenOpts)
			if err != nil {
				logger.Error("interactive authentication failed",
					logging.String("event", eventAuthLoginTokenFailure),
					logging.Int("attempt", attempt),
					logging.Error(err),
				)
				return nil, fmt.Errorf("authentication failed: %w", err)
			}

			record = newRecord

			cred, err = s.factory.Create(authinfra.CredentialOptions{
				Type:                 cfg.Auth.Type,
				ClientID:             cfg.Auth.ClientID,
				TenantID:             cfg.Auth.TenantID,
				AuthenticationRecord: newRecord,
			})
			if err != nil {
				logger.Error("failed to reload credential after authentication", logging.Error(err))
				return nil, fmt.Errorf("failed to reload credential: %w", err)
			}
		}

		logger.Debug("requesting access token",
			logging.String("event", eventAuthLoginTokenAttempt),
			logging.Int("attempt", attempt),
			logging.Int("max_attempts", maxAttempts),
		)

		token, err = cred.GetToken(ctx, *tokenOpts)
		if err != nil {
			if !isAuthRequired(err) {
				logger.Error("failed to retrieve token",
					logging.String("event", eventAuthLoginTokenFailure),
					logging.Int("attempt", attempt),
					logging.Error(err),
				)
				return nil, fmt.Errorf("failed to retrieve token: %w", err)
			}

			logger.Warn("token request indicates authentication required; clearing record and retrying",
				logging.String("event", eventAuthLoginTokenFailure),
				logging.Int("attempt", attempt),
				logging.Error(err),
			)

			record = azidentity.AuthenticationRecord{}
			continue
		}

		logger.Info("token retrieved successfully",
			logging.String("event", eventAuthLoginTokenSuccess),
			logging.Int("attempt", attempt),
			logging.String("expires_on", token.ExpiresOn.String()),
		)
		break
	}

	if token == (azcore.AccessToken{}) {
		logger.Error("token is empty after login attempts",
			logging.String("event", eventAuthLoginEmptyToken),
		)
		return nil, errors.New("empty token")
	}

	if err := s.cache.SetProfile(ctx, profileName, record); err != nil {
		logger.Warn("failed to save authentication record",
			logging.String("event", eventAuthLoginRecordSaveErr),
			logging.Error(err),
		)
	} else {
		logger.Info("authentication record saved",
			logging.String("event", eventAuthLoginRecordSaved),
			logging.Bool("record_saved", record != (azidentity.AuthenticationRecord{})),
		)
	}

	logger.Info("login successful",
		logging.String("event", eventAuthLoginSuccess),
		logging.Bool("record_saved", record != (azidentity.AuthenticationRecord{})),
	)

	return &authdomain.LoginResult{
		AccessToken: token.Token,
		RecordSaved: (record != azidentity.AuthenticationRecord{}),
	}, nil
}

// Logout removes the cached authentication record for the given profile.
// If force is true, logout proceeds even if no record exists.
func (s *AuthService) Logout(ctx context.Context, profileName string, force bool) error {
	correlationID := util.CorrelationIDFromContext(ctx)

	logger := s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("profile", profileName),
		logging.String("event", eventAuthLogoutStart),
	)

	logger.Info("starting logout flow", logging.Bool("force", force))

	record, err := s.cache.GetProfile(ctx, profileName)
	if err != nil {
		if force {
			logger.Warn("failed to load cached record, but force=true; continuing",
				logging.String("event", eventAuthLogoutRecordLoad),
				logging.Error(err),
			)
		} else {
			logger.Error("failed to load cached record",
				logging.String("event", eventAuthLogoutRecordLoad),
				logging.Error(err),
			)
			return fmt.Errorf("failed to load cached authentication record: %w", err)
		}
	} else {
		logger.Debug("loaded cached authentication record for logout",
			logging.String("event", eventAuthLogoutRecordLoad),
			logging.Bool("has_record", record != (azidentity.AuthenticationRecord{})),
		)
	}

	if record == (azidentity.AuthenticationRecord{}) && !force {
		logger.Info("no authentication record found; nothing to do",
			logging.String("event", eventAuthLogoutRecordMissing),
		)
		return nil
	}

	logger.Info("deleting cached authentication record",
		logging.String("event", eventAuthLogoutDelete),
	)

	if err := s.cache.DeleteProfile(ctx, profileName); err != nil {
		logger.Error("failed to delete authentication record",
			logging.String("event", eventAuthLogoutDeleteError),
			logging.Error(err),
		)
		return fmt.Errorf("failed to delete authentication record: %w", err)
	}

	logger.Info("authentication record deleted",
		logging.String("event", eventAuthLogoutDelete),
	)

	logger.Info("logout successful",
		logging.String("event", eventAuthLogoutSuccess),
	)

	return nil
}
