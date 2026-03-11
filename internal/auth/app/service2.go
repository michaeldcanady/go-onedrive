package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	domainaccount "github.com/michaeldcanady/go-onedrive/internal/account/domain"
	domainauth "github.com/michaeldcanady/go-onedrive/internal/auth/domain"
	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainconfig "github.com/michaeldcanady/go-onedrive/internal/config/domain"
	logger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
	"github.com/michaeldcanady/go-onedrive/pkg/cache"
	pkgcache "github.com/michaeldcanady/go-onedrive/pkg/cache"
)

var (
	ErrSilentRequiresAccount = errors.New("silent token acquisition requires a cached account")
	ErrMaxAuthAttempts       = errors.New("max authentication attempts reached")
	ErrEmptyToken            = errors.New("empty token returned from provider")
)

const (
	eventAuthSilentStart = "domainauth.silent.start"
	eventAuthLoginStart  = "domainauth.login.start"
	eventAuthLogoutStart = "domainauth.logout.start"
)

type Service2 struct {
	cache   pkgcache.Cache[domainauth.AccessToken]
	config  domainconfig.ConfigService
	state   domainstate.Service
	log     logger.Logger
	factory domainauth.CredentialFactory
	account domainaccount.Service
}

func NewService(
	cache pkgcache.Cache[domainauth.AccessToken],
	config domainconfig.ConfigService,
	state domainstate.Service,
	l logger.Logger,
	factory domainauth.CredentialFactory,
	account domainaccount.Service,
) *Service2 {
	return &Service2{
		cache:   cache,
		config:  config,
		state:   state,
		log:     l,
		factory: factory,
		account: account,
	}
}

func (s *Service2) buildLogger(ctx context.Context) logger.Logger {
	correlationID := util.CorrelationIDFromContext(ctx)
	return s.log.WithContext(ctx).With(
		logger.String("correlation_id", correlationID),
	)
}

func (s *Service2) resolveProfile(profileName string) (string, error) {
	if profileName != "" {
		return profileName, nil
	}
	return s.state.Get(domainstate.KeyProfile)
}

func (s *Service2) loadProfileConfig(
	ctx context.Context,
	log logger.Logger,
	profileName string,
) (string, domainconfig.Configuration, error) {

	profile, err := s.resolveProfile(profileName)
	if err != nil {
		log.Warn("failed to resolve profile", logger.Error(err))
		return "", domainconfig.Configuration{}, fmt.Errorf("failed to resolve profile: %w", err)
	}

	cfg, err := s.config.GetConfiguration(ctx, profile)
	if err != nil {
		log.Warn("failed to load configuration",
			logger.String("profile", profile),
			logger.Error(err),
		)
		return "", domainconfig.Configuration{}, fmt.Errorf("failed to load configuration: %w", err)
	}

	return profile, cfg, nil
}

func (s *Service2) buildCredentialProvider(
	log logger.Logger,
	cfg domainconfig.Configuration,
	account domainaccount.Account,
) (domainauth.CredentialProvider, error) {

	opts := &domainauth.Options{
		Method:   cfg.Auth.Type,
		TenantID: cfg.Auth.TenantID,
		ClientID: cfg.Auth.ClientID,
	}

	provider, err := s.factory.Credential(account, opts)
	if err != nil {
		log.Warn("failed to build credential provider", logger.Error(err))
		return nil, fmt.Errorf("failed to build credential provider: %w", err)
	}

	log.Info("credential provider initialized")
	return provider, nil
}

func (s *Service2) loadAccountOrEmpty(ctx context.Context, log logger.Logger) domainaccount.Account {
	acc, err := s.account.Get(ctx)
	if err != nil {
		log.Debug("no cached account found", logger.Error(err))
		return domainaccount.Account{}
	}
	return acc
}

func (s *Service2) buildTokenRequestOptions(opts domainauth.LoginOptions) *policy.TokenRequestOptions {
	return &policy.TokenRequestOptions{
		Scopes:    opts.Scopes,
		EnableCAE: opts.EnableCAE,
	}
}

func (s *Service2) needsInteractiveAuth(opts domainauth.LoginOptions, record domainaccount.Account) bool {
	return opts.Force || record == (domainaccount.Account{})
}

func (s *Service2) getCachedToken(ctx context.Context, account domainaccount.Account) (domainauth.AccessToken, error) {
	token, err := s.cache.Get(ctx, account.HomeAccountID)
	if err != nil {
		if !errors.Is(err, cache.ErrKeyNotFound) {
			return domainauth.AccessToken{}, err
		}
	}

	if time.Now().Before(token.ExpiresOn) && token.Token != "" {
		return token, nil
	}

	return domainauth.AccessToken{}, nil
}

func (s *Service2) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	log := s.buildLogger(ctx).With(logger.String("event", eventAuthSilentStart))
	log.Info("starting silent token acquisition")

	profile, cfg, err := s.loadProfileConfig(ctx, log, "")
	if err != nil {
		return azcore.AccessToken{}, err
	}
	log = log.With(logger.String("profile", profile))

	account := s.loadAccountOrEmpty(ctx, log)
	if account == (domainaccount.Account{}) {
		log.Warn("silent auth requires cached account")
		return azcore.AccessToken{}, ErrSilentRequiresAccount
	}

	token, err := s.getCachedToken(ctx, account)
	if err != nil {
		return azcore.AccessToken{}, err
	}

	if token != (domainauth.AccessToken{}) {
		return azcore.AccessToken{
			Token:     token.Token,
			ExpiresOn: token.ExpiresOn,
			RefreshOn: token.RefreshOn,
		}, nil
	}

	provider, err := s.buildCredentialProvider(log, cfg, account)
	if err != nil {
		return azcore.AccessToken{}, err
	}

	azToken, err := provider.GetToken(ctx, options)
	if err != nil {
		log.Error("silent token acquisition failed", logger.Error(err))
		return azcore.AccessToken{}, err
	}

	token = domainauth.AccessToken{
		Token:     azToken.Token,
		ExpiresOn: azToken.ExpiresOn,
		RefreshOn: azToken.ExpiresOn,
	}

	log.Info("silent token acquisition successful",
		logger.String("expires_on", token.ExpiresOn.String()),
	)

	return azToken, nil
}

func (s *Service2) Login(ctx context.Context, profileName string, opts domainauth.LoginOptions) (*domainauth.LoginResult, error) {
	log := s.buildLogger(ctx).With(logger.String("event", eventAuthLoginStart))
	log.Info("starting login flow")

	profile, cfg, err := s.loadProfileConfig(ctx, log, profileName)
	if err != nil {
		return nil, err
	}
	log = log.With(logger.String("profile", profile))

	account := s.loadAccountOrEmpty(ctx, log)

	token, err := s.getCachedToken(ctx, account)
	if err != nil {
		return nil, err
	}

	if token != (domainauth.AccessToken{}) {
		return &domainauth.LoginResult{
			AccessToken: token.Token,
			RecordSaved: true,
			ExpiresOn:   token.ExpiresOn,
			Username:    account.Username,
			Profile:     profile,
		}, nil
	}

	provider, err := s.buildCredentialProvider(log, cfg, account)
	if err != nil {
		return nil, err
	}

	tokenOpts := s.buildTokenRequestOptions(opts)

	token, record, err := s.performLoginFlow(ctx, log, provider, tokenOpts, opts, account)
	if err != nil {
		return nil, err
	}

	if err := s.validateLoginToken(log, token); err != nil {
		return nil, err
	}

	if err := s.account.Put(ctx, record); err != nil {
		log.Warn("failed to cache account", logger.Error(err))
	}

	if err := s.cache.Set(ctx, record.HomeAccountID, token); err != nil {
		log.Warn("failed to cache token", logger.Error(err))
	}

	log.Info("login successful")

	return &domainauth.LoginResult{
		AccessToken: token.Token,
		ExpiresOn:   token.ExpiresOn,
		Username:    record.Username,
		Profile:     profile,
		RecordSaved: record != (domainaccount.Account{}),
	}, nil
}

func (s *Service2) performLoginFlow(
	ctx context.Context,
	log logger.Logger,
	provider domainauth.CredentialProvider,
	tokenOpts *policy.TokenRequestOptions,
	opts domainauth.LoginOptions,
	record domainaccount.Account,
) (domainauth.AccessToken, domainaccount.Account, error) {

	var (
		token       domainauth.AccessToken
		maxAttempts = int32(3)
		err         error
	)

	attempt := int32(0)
	for {
		if ctx.Err() != nil {
			return domainauth.AccessToken{}, domainaccount.Account{}, ctx.Err()
		}
		attempt++

		if attempt >= maxAttempts {
			return domainauth.AccessToken{}, domainaccount.Account{}, ErrMaxAuthAttempts
		}

		log.Debug("login attempt",
			logger.Int("attempt", int(attempt)),
			logger.Int("max_attempts", int(maxAttempts)),
		)

		if s.needsInteractiveAuth(opts, record) {
			record, err = s.performInteractiveAuthAttempt(ctx, log, provider, tokenOpts)
			if err != nil {
				continue
			}
		}

		azToken, err := provider.GetToken(ctx, *tokenOpts)

		token = domainauth.AccessToken{
			Token:     azToken.Token,
			ExpiresOn: azToken.ExpiresOn,
			RefreshOn: azToken.RefreshOn,
		}

		if err == nil {
			break
		}

		if !isAuthRequired(err) {
			log.Error("token retrieval failed", logger.Error(err))
			return domainauth.AccessToken{}, domainaccount.Account{}, fmt.Errorf("token retrieval failed: %w", err)
		}

		log.Warn("authentication required; clearing record and retrying", logger.Error(err))
		record = domainaccount.Account{}
	}

	return token, record, nil
}

func (s *Service2) performInteractiveAuthAttempt(
	ctx context.Context,
	log logger.Logger,
	provider domainauth.CredentialProvider,
	tokenOpts *policy.TokenRequestOptions,
) (domainaccount.Account, error) {

	log.Info("performing interactive authentication")

	newRecord, err := provider.Authenticate(ctx, tokenOpts)
	if err != nil {
		log.Warn("interactive authentication failed", logger.Error(err))
		return domainaccount.Account{}, fmt.Errorf("interactive authentication failed: %w", err)
	}

	account := domainaccount.AccountFromMSAuthRecord(newRecord)

	log.Info("interactive authentication successful",
		logger.String("username", account.Username),
	)

	return account, nil
}

func (s *Service2) validateLoginToken(log logger.Logger, token domainauth.AccessToken) error {
	if token.Token == "" || token.ExpiresOn.IsZero() {
		log.Error("empty or invalid token returned")
		return ErrEmptyToken
	}
	return nil
}

func (s *Service2) Logout(ctx context.Context, profileName string, force bool) error {
	log := s.buildLogger(ctx).With(
		logger.String("profile", profileName),
	)

	log.Info("starting logout flow",
		logger.String("event", eventAuthLogoutStart),
	)

	account := s.loadAccountOrEmpty(ctx, log)
	if account == (domainaccount.Account{}) {
		log.Info("no cached account; nothing to do")
		return nil
	}

	if err := s.account.Delete(ctx); err != nil {
		log.Warn("failed to delete account", logger.Error(err))
		return fmt.Errorf("failed to delete account: %w", err)
	}

	if err := s.cache.Delete(ctx, account.HomeAccountID); err != nil {
		log.Warn("failed to delete token", logger.Error(err))
		return fmt.Errorf("failed to delete token: %w", err)
	}

	log.Info("logout successful")
	return nil
}
