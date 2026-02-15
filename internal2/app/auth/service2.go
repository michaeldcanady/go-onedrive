package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	authdomain "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	domainconfig "github.com/michaeldcanady/go-onedrive/internal2/domain/config"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/state"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type Service2 struct {
	cache   *abstractions.Cache2
	config  domainconfig.ConfigService
	state   state.Service
	logger  logging.Logger
	factory auth.CredentialFactory
	account cache.CacheService
}

func NewService2(
	cache *abstractions.Cache2,
	config domainconfig.ConfigService,
	state state.Service,
	logger logging.Logger,
	factory auth.CredentialFactory,
	account cache.CacheService,
) *Service2 {
	return &Service2{
		config:  config,
		state:   state,
		logger:  logger,
		factory: factory,
		account: account,
	}
}

// ---------- shared helpers ----------

func (s *Service2) buildLogger(ctx context.Context) logging.Logger {
	correlationID := util.CorrelationIDFromContext(ctx)
	return s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
	)
}

func (s *Service2) loadConfiguration(ctx context.Context, logger logging.Logger, profile string) (config.Configuration3, error) {
	logger.Debug("loading configuration",
		logging.String("profile", profile),
	)

	cfg, err := s.config.GetConfiguration(ctx, profile)
	if err != nil {
		logger.Warn("failed to retrieve config",
			logging.Error(err),
		)
		return config.Configuration3{}, err
	}

	return cfg, nil
}

func (s *Service2) getCredential(
	cfg config.Configuration3,
	account azidentity.AuthenticationRecord,
) (auth.CredentialProvider, error) {
	opts := &auth.Options{
		Method:   auth.ParseMethod(cfg.Auth.Type),
		TenantID: cfg.Auth.TenantID,
		ClientID: cfg.Auth.ClientID,
	}

	return s.factory.Credential(account, opts)
}

func (s *Service2) buildCredentialProvider(
	ctx context.Context,
	logger logging.Logger,
	cfg config.Configuration3,
	account azidentity.AuthenticationRecord,
) (auth.CredentialProvider, error) {
	logger.Debug("retrieving credential provider")
	provider, err := s.getCredential(cfg, account)
	if err != nil {
		logger.Warn("failed to retrieve credential",
			logging.Error(err),
		)
		return nil, err
	}
	return provider, nil
}

// ---------- GetToken ----------

func (s *Service2) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	logger := s.buildLogger(ctx)
	s.logSilentStart(logger, options)

	account, err := s.loadCachedAccountForSilent(ctx, logger, "default")
	if err != nil {
		logger.Warn("failed to load cached account",
			logging.String("profile", "default"),
			logging.Error(err),
		)

		return azcore.AccessToken{}, err
	}

	cfg, err := s.loadConfiguration(ctx, logger, "default")
	if err != nil {
		logger.Warn("failed to load configuration",
			logging.String("profile", "default"),
			logging.Error(err),
		)

		return azcore.AccessToken{}, err
	}

	provider, err := s.buildCredentialProvider(ctx, logger, cfg, account)
	if err != nil {
		return azcore.AccessToken{}, err
	}

	return provider.GetToken(ctx, options)
}

func (s *Service2) logSilentStart(logger logging.Logger, options policy.TokenRequestOptions) {
	logger.Info("starting silent token acquisition",
		logging.String("event", eventAuthSilentStart),
		logging.String("scopes", fmt.Sprintf("%v", options.Scopes)),
	)
}

func (s *Service2) loadCachedAccountForSilent(ctx context.Context, logger logging.Logger, profile string) (azidentity.AuthenticationRecord, error) {
	logger.Debug("loading cached authentication record for silent token acquisition")

	// TODO: With new cache system how to get cached profile?
	return s.account.GetProfile(ctx, profile)
}

// ---------- Login ----------

func (s *Service2) Login(ctx context.Context, profileName string, opts authdomain.LoginOptions) (*authdomain.LoginResult, error) {
	logger := s.buildLogger(ctx)
	s.logLoginStart(logger, opts)

	account := s.loadCachedAccountForLogin(ctx, logger)

	cfg, err := s.loadConfiguration(ctx, logger, "default")
	if err != nil {
		return nil, err
	}

	provider, err := s.buildCredentialProvider(ctx, logger, cfg, account)
	if err != nil {
		return nil, err
	}

	tokenOpts := s.buildTokenRequestOptions(opts)

	token, record, err := s.performLoginFlow(ctx, logger, provider, cfg, tokenOpts, opts)
	if err != nil {
		return nil, err
	}

	if err := s.validateLoginToken(logger, token); err != nil {
		return nil, err
	}

	s.logLoginSuccess(logger, record)

	// TODO: cache account
	// TODO: cache token

	return &authdomain.LoginResult{
		AccessToken: token.Token,
		RecordSaved: (record != azidentity.AuthenticationRecord{}),
	}, nil
}

func (s *Service2) logLoginStart(logger logging.Logger, opts authdomain.LoginOptions) {
	logger.Info("starting login flow",
		logging.Bool("force", opts.Force),
		logging.Bool("enable_cae", opts.EnableCAE),
		logging.String("scopes", fmt.Sprintf("%v", opts.Scopes)),
	)
}

func (s *Service2) loadCachedAccountForLogin(ctx context.Context, logger logging.Logger) azidentity.AuthenticationRecord {
	// TODO: With new cache system how to get cached profile?
	var account azidentity.AuthenticationRecord
	logger.Debug("loading cached authentication record for login")
	return account
}

func (s *Service2) buildTokenRequestOptions(opts authdomain.LoginOptions) *policy.TokenRequestOptions {
	return &policy.TokenRequestOptions{
		Scopes:    opts.Scopes,
		EnableCAE: opts.EnableCAE,
	}
}

func (s *Service2) performLoginFlow(
	ctx context.Context,
	logger logging.Logger,
	provider auth.CredentialProvider,
	cfg config.Configuration3,
	tokenOpts *policy.TokenRequestOptions,
	opts authdomain.LoginOptions,
) (azcore.AccessToken, azidentity.AuthenticationRecord, error) {
	var (
		token       azcore.AccessToken
		record      azidentity.AuthenticationRecord
		maxAttempts = int32(3)
		err         error
	)

	for attempt := int32(1); attempt <= maxAttempts; attempt++ {
		if s.needsInteractiveAuth(opts, record) {
			record, provider, err = s.performInteractiveAuthAttempt(ctx, logger, provider, cfg, tokenOpts, attempt, maxAttempts)
			if err != nil {
				return azcore.AccessToken{}, azidentity.AuthenticationRecord{}, err
			}
		}

		token, err = s.requestAccessToken(ctx, logger, provider, tokenOpts, attempt, maxAttempts)
		if err == nil {
			break
		}

		if !isAuthRequired(err) {
			return azcore.AccessToken{}, azidentity.AuthenticationRecord{}, s.handleNonAuthError(logger, err, attempt)
		}

		s.logAuthRequiredRetry(logger, err, attempt)
		record = azidentity.AuthenticationRecord{}
	}

	return token, record, nil
}

func (s *Service2) needsInteractiveAuth(opts authdomain.LoginOptions, record azidentity.AuthenticationRecord) bool {
	return opts.Force || record == (azidentity.AuthenticationRecord{})
}

func (s *Service2) performInteractiveAuthAttempt(
	ctx context.Context,
	logger logging.Logger,
	provider auth.CredentialProvider,
	cfg config.Configuration3,
	tokenOpts *policy.TokenRequestOptions,
	attempt, maxAttempts int32,
) (azidentity.AuthenticationRecord, auth.CredentialProvider, error) {
	s.logInteractiveAuthStart(logger, attempt, maxAttempts)

	newRecord, err := provider.Authenticate(ctx, tokenOpts)
	if err != nil {
		s.logInteractiveAuthFailure(logger, err, attempt)
		return azidentity.AuthenticationRecord{}, nil, fmt.Errorf("authentication failed: %w", err)
	}

	provider, err = s.reloadCredentialAfterAuth(ctx, logger, cfg, newRecord)
	if err != nil {
		return azidentity.AuthenticationRecord{}, nil, err
	}

	return newRecord, provider, nil
}

func (s *Service2) logInteractiveAuthStart(logger logging.Logger, attempt, maxAttempts int32) {
	logger.Info("performing interactive authentication",
		logging.String("event", eventAuthLoginInteractive),
		logging.Int("attempt", attempt),
		logging.Int("max_attempts", maxAttempts),
	)
}

func (s *Service2) logInteractiveAuthFailure(logger logging.Logger, err error, attempt int32) {
	logger.Error("interactive authentication failed",
		logging.String("event", eventAuthLoginTokenFailure),
		logging.Int("attempt", attempt),
		logging.Error(err),
	)
}

func (s *Service2) reloadCredentialAfterAuth(
	ctx context.Context,
	logger logging.Logger,
	cfg config.Configuration3,
	record azidentity.AuthenticationRecord,
) (auth.CredentialProvider, error) {
	provider, err := s.getCredential(cfg, record)
	if err != nil {
		logger.Error("failed to reload credential after authentication", logging.Error(err))
		return nil, fmt.Errorf("failed to reload credential: %w", err)
	}
	return provider, nil
}

func (s *Service2) requestAccessToken(
	ctx context.Context,
	logger logging.Logger,
	provider auth.CredentialProvider,
	tokenOpts *policy.TokenRequestOptions,
	attempt, maxAttempts int32,
) (azcore.AccessToken, error) {
	logger.Debug("requesting access token",
		logging.String("event", eventAuthLoginTokenAttempt),
		logging.Int("attempt", attempt),
		logging.Int("max_attempts", maxAttempts),
	)

	token, err := provider.GetToken(ctx, *tokenOpts)
	if err != nil {
		return azcore.AccessToken{}, err
	}

	logger.Info("token retrieved successfully",
		logging.String("event", eventAuthLoginTokenSuccess),
		logging.Int("attempt", attempt),
		logging.String("expires_on", token.ExpiresOn.String()),
	)

	return token, nil
}

func (s *Service2) handleNonAuthError(logger logging.Logger, err error, attempt int32) error {
	logger.Error("failed to retrieve token",
		logging.String("event", eventAuthLoginTokenFailure),
		logging.Int("attempt", attempt),
		logging.Error(err),
	)
	return fmt.Errorf("failed to retrieve token: %w", err)
}

func (s *Service2) logAuthRequiredRetry(logger logging.Logger, err error, attempt int32) {
	logger.Warn("token request indicates authentication required; clearing record and retrying",
		logging.String("event", eventAuthLoginTokenFailure),
		logging.Int("attempt", attempt),
		logging.Error(err),
	)
}

func (s *Service2) validateLoginToken(logger logging.Logger, token azcore.AccessToken) error {
	if token == (azcore.AccessToken{}) {
		logger.Error("token is empty after login attempts",
			logging.String("event", eventAuthLoginEmptyToken),
		)
		return errors.New("empty token")
	}
	return nil
}

func (s *Service2) logLoginSuccess(logger logging.Logger, record azidentity.AuthenticationRecord) {
	logger.Info("login successful",
		logging.String("event", eventAuthLoginSuccess),
		logging.Bool("record_saved", record != (azidentity.AuthenticationRecord{})),
	)
}
