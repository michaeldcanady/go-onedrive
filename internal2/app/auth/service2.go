package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	accountdomain "github.com/michaeldcanady/go-onedrive/internal2/domain/account"
	authdomain "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	domainconfig "github.com/michaeldcanady/go-onedrive/internal2/domain/config"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/state"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/core"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/config"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

var (
	ErrSilentRequiresAccount = errors.New("silent token acquisition requires a cached account")
	ErrMaxAuthAttempts       = errors.New("max authentication attempts reached")
	ErrEmptyToken            = errors.New("empty token returned from provider")
)

type Service2 struct {
	cache   *abstractions.Cache2
	config  domainconfig.ConfigService
	state   state.Service
	logger  logging.Logger
	factory authdomain.CredentialFactory
	account accountdomain.Service
}

func NewService2(
	cache *abstractions.Cache2,
	config domainconfig.ConfigService,
	state state.Service,
	logger logging.Logger,
	factory authdomain.CredentialFactory,
	account accountdomain.Service,
) *Service2 {
	return &Service2{
		cache:   cache,
		config:  config,
		state:   state,
		logger:  logger,
		factory: factory,
		account: account,
	}
}

func (s *Service2) buildLogger(ctx context.Context) logging.Logger {
	correlationID := util.CorrelationIDFromContext(ctx)
	return s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
	)
}

func (s *Service2) resolveProfile(profileName string) (string, error) {
	if profileName != "" {
		return profileName, nil
	}
	return s.state.GetCurrentProfile()
}

func (s *Service2) loadProfileConfig(
	ctx context.Context,
	logger logging.Logger,
	profileName string,
) (string, config.Configuration3, error) {

	profile, err := s.resolveProfile(profileName)
	if err != nil {
		logger.Warn("failed to resolve profile", logging.Error(err))
		return "", config.Configuration3{}, fmt.Errorf("failed to resolve profile: %w", err)
	}

	cfg, err := s.config.GetConfiguration(ctx, profile)
	if err != nil {
		logger.Warn("failed to load configuration",
			logging.String("profile", profile),
			logging.Error(err),
		)
		return "", config.Configuration3{}, fmt.Errorf("failed to load configuration: %w", err)
	}

	return profile, cfg, nil
}

func (s *Service2) buildCredentialProvider(
	logger logging.Logger,
	cfg config.Configuration3,
	account accountdomain.Account,
) (authdomain.CredentialProvider, error) {

	opts := &authdomain.Options{
		Method:   authdomain.ParseMethod(cfg.Auth.Type),
		TenantID: cfg.Auth.TenantID,
		ClientID: cfg.Auth.ClientID,
	}

	provider, err := s.factory.Credential(account, opts)
	if err != nil {
		logger.Warn("failed to build credential provider", logging.Error(err))
		return nil, fmt.Errorf("failed to build credential provider: %w", err)
	}

	logger.Info("credential provider initialized")
	return provider, nil
}

func (s *Service2) loadAccountOrEmpty(ctx context.Context, logger logging.Logger) accountdomain.Account {
	acc, err := s.account.Get(ctx)
	if err != nil {
		logger.Debug("no cached account found", logging.Error(err))
		return accountdomain.Account{}
	}
	return acc
}

func (s *Service2) buildTokenRequestOptions(opts authdomain.LoginOptions) *policy.TokenRequestOptions {
	return &policy.TokenRequestOptions{
		Scopes:    opts.Scopes,
		EnableCAE: opts.EnableCAE,
	}
}

func (s *Service2) needsInteractiveAuth(opts authdomain.LoginOptions, record accountdomain.Account) bool {
	return opts.Force || record == (accountdomain.Account{})
}

func (s *Service2) getCachedToken(ctx context.Context, account accountdomain.Account) (authdomain.AccessToken, error) {
	var token authdomain.AccessToken

	if err := s.cache.Get(ctx, func() ([]byte, error) { return json.Marshal(account.HomeAccountID) }, func(b []byte) error { return json.Unmarshal(b, &token) }); err != nil {
		if !errors.Is(err, core.ErrKeyNotFound) {
			return authdomain.AccessToken{}, err
		}
	}

	if time.Now().Before(token.ExpiresOn) && token.Token != "" {
		return token, nil
	}

	return authdomain.AccessToken{}, nil
}

func (s *Service2) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	logger := s.buildLogger(ctx).With(logging.String("event", eventAuthSilentStart))
	logger.Info("starting silent token acquisition")

	profile, cfg, err := s.loadProfileConfig(ctx, logger, "")
	if err != nil {
		return azcore.AccessToken{}, err
	}
	logger = logger.With(logging.String("profile", profile))

	account := s.loadAccountOrEmpty(ctx, logger)
	if account == (accountdomain.Account{}) {
		logger.Warn("silent auth requires cached account")
		return azcore.AccessToken{}, ErrSilentRequiresAccount
	}

	token, err := s.getCachedToken(ctx, account)
	if err != nil {
		return azcore.AccessToken{}, err
	}

	if token != (authdomain.AccessToken{}) {
		return azcore.AccessToken{
			Token:     token.Token,
			ExpiresOn: token.ExpiresOn,
			RefreshOn: token.RefreshOn,
		}, nil
	}

	provider, err := s.buildCredentialProvider(logger, cfg, account)
	if err != nil {
		return azcore.AccessToken{}, err
	}

	azToken, err := provider.GetToken(ctx, options)
	if err != nil {
		logger.Error("silent token acquisition failed", logging.Error(err))
		return azcore.AccessToken{}, err
	}

	token = authdomain.AccessToken{
		Token:     azToken.Token,
		ExpiresOn: azToken.ExpiresOn,
		RefreshOn: azToken.ExpiresOn,
	}

	logger.Info("silent token acquisition successful",
		logging.String("expires_on", token.ExpiresOn.String()),
	)

	return azToken, nil
}

func (s *Service2) Login(ctx context.Context, profileName string, opts authdomain.LoginOptions) (*authdomain.LoginResult, error) {
	logger := s.buildLogger(ctx).With(logging.String("event", eventAuthLoginStart))
	logger.Info("starting login flow")

	profile, cfg, err := s.loadProfileConfig(ctx, logger, profileName)
	if err != nil {
		return nil, err
	}
	logger = logger.With(logging.String("profile", profile))

	account := s.loadAccountOrEmpty(ctx, logger)

	token, err := s.getCachedToken(ctx, account)
	if err != nil {
		return nil, err
	}

	if token != (authdomain.AccessToken{}) {
		return &authdomain.LoginResult{
			AccessToken: token.Token,
			RecordSaved: true,
			ExpiresOn:   token.ExpiresOn,
			Username:    account.Username,
			Profile:     profile,
		}, nil
	}

	provider, err := s.buildCredentialProvider(logger, cfg, account)
	if err != nil {
		return nil, err
	}

	tokenOpts := s.buildTokenRequestOptions(opts)

	token, record, err := s.performLoginFlow(ctx, logger, provider, tokenOpts, opts, account)
	if err != nil {
		return nil, err
	}

	if err := s.validateLoginToken(logger, token); err != nil {
		return nil, err
	}

	if err := s.account.Put(ctx, record); err != nil {
		logger.Warn("failed to cache account", logging.Error(err))
	}

	s.cache.Set(ctx, func() ([]byte, error) { return json.Marshal(record.HomeAccountID) }, func() ([]byte, error) { return json.Marshal(token) })

	logger.Info("login successful")

	return &authdomain.LoginResult{
		AccessToken: token.Token,
		ExpiresOn:   token.ExpiresOn,
		Username:    record.Username,
		Profile:     profile,
		RecordSaved: record != (accountdomain.Account{}),
	}, nil
}

func (s *Service2) performLoginFlow(
	ctx context.Context,
	logger logging.Logger,
	provider authdomain.CredentialProvider,
	tokenOpts *policy.TokenRequestOptions,
	opts authdomain.LoginOptions,
	record accountdomain.Account,
) (authdomain.AccessToken, accountdomain.Account, error) {

	var (
		token       authdomain.AccessToken
		maxAttempts = int32(3)
		err         error
	)

	attempt := int32(0)
	for {
		if ctx.Err() != nil {
			return authdomain.AccessToken{}, accountdomain.Account{}, ctx.Err()
		}
		attempt++

		if attempt >= maxAttempts {
			return authdomain.AccessToken{}, accountdomain.Account{}, ErrMaxAuthAttempts
		}

		logger.Debug("login attempt",
			logging.Int("attempt", int(attempt)),
			logging.Int("max_attempts", int(maxAttempts)),
		)

		if s.needsInteractiveAuth(opts, record) {
			record, err = s.performInteractiveAuthAttempt(ctx, logger, provider, tokenOpts)
			if err != nil {
				continue
			}
		}

		azToken, err := provider.GetToken(ctx, *tokenOpts)

		token = authdomain.AccessToken{
			Token:     azToken.Token,
			ExpiresOn: azToken.ExpiresOn,
			RefreshOn: azToken.RefreshOn,
		}

		if err == nil {
			break
		}

		if !isAuthRequired(err) {
			logger.Error("token retrieval failed", logging.Error(err))
			return authdomain.AccessToken{}, accountdomain.Account{}, fmt.Errorf("token retrieval failed: %w", err)
		}

		logger.Warn("authentication required; clearing record and retrying", logging.Error(err))
		record = accountdomain.Account{}
	}

	return token, record, nil
}

func (s *Service2) performInteractiveAuthAttempt(
	ctx context.Context,
	logger logging.Logger,
	provider authdomain.CredentialProvider,
	tokenOpts *policy.TokenRequestOptions,
) (accountdomain.Account, error) {

	logger.Info("performing interactive authentication")

	newRecord, err := provider.Authenticate(ctx, tokenOpts)
	if err != nil {
		logger.Warn("interactive authentication failed", logging.Error(err))
		return accountdomain.Account{}, fmt.Errorf("interactive authentication failed: %w", err)
	}

	account := accountdomain.AccountFromMSAuthRecord(newRecord)

	logger.Info("interactive authentication successful",
		logging.String("username", account.Username),
	)

	return account, nil
}

func (s *Service2) validateLoginToken(logger logging.Logger, token authdomain.AccessToken) error {
	if token.Token == "" || token.ExpiresOn.IsZero() {
		logger.Error("empty or invalid token returned")
		return ErrEmptyToken
	}
	return nil
}

func (s *Service2) Logout(ctx context.Context, profileName string, force bool) error {
	logger := s.buildLogger(ctx).With(
		logging.String("profile", profileName),
	)

	logger.Info("starting logout flow",
		logging.String("event", eventAuthLogoutStart),
	)

	account := s.loadAccountOrEmpty(ctx, logger)
	if account == (accountdomain.Account{}) {
		logger.Info("no cached account; nothing to do")
		return nil
	}

	if err := s.account.Delete(ctx); err != nil {
		logger.Warn("failed to delete account", logging.Error(err))
		return fmt.Errorf("failed to delete account: %w", err)
	}

	if err := s.cache.Delete(ctx, func() ([]byte, error) { return json.Marshal(account.HomeAccountID) }); err != nil {
		logger.Warn("failed to delete token", logging.Error(err))
		return fmt.Errorf("failed to delete token: %w", err)
	}

	logger.Info("logout successful")
	return nil
}
