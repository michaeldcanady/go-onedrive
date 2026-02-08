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
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Info("starting silent token acquisition",
		logging.String("event", "get_token_start"),
		logging.String("correlation_id", cid),
		logging.Strings("scopes", options.Scopes),
	)

	profileName, err := s.state.GetCurrentProfile()
	if err != nil {
		s.logger.Error("failed to resolve current profile",
			logging.String("event", "resolve_profile"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return azcore.AccessToken{}, err
	}

	s.logger.Debug("loading cached authentication record",
		logging.String("event", "load_cached_record"),
		logging.String("profile", profileName),
		logging.String("correlation_id", cid),
	)

	record, err := s.cache.GetProfile(ctx, profileName)
	if err != nil {
		s.logger.Warn("failed to load cached authentication record",
			logging.String("event", "load_cached_record"),
			logging.Error(err),
			logging.String("profile", profileName),
			logging.String("correlation_id", cid),
		)
		return azcore.AccessToken{}, fmt.Errorf("failed to load cached authentication record: %w", err)
	}

	s.logger.Debug("loading profile configuration",
		logging.String("event", "load_config"),
		logging.String("profile", profileName),
		logging.String("correlation_id", cid),
	)

	cfg, err := s.config.GetConfiguration(ctx, profileName)
	if err != nil {
		s.logger.Error("failed to load configuration",
			logging.String("event", "load_config"),
			logging.Error(err),
			logging.String("profile", profileName),
			logging.String("correlation_id", cid),
		)
		return azcore.AccessToken{}, fmt.Errorf("failed to load configuration: %w", err)
	}

	credOpts := authinfra.CredentialOptions{
		Type:                 cfg.Auth.Type,
		ClientID:             cfg.Auth.ClientID,
		TenantID:             cfg.Auth.TenantID,
		AuthenticationRecord: record,
	}

	s.logger.Debug("creating credential",
		logging.String("event", "create_credential"),
		logging.String("auth_type", string(cfg.Auth.Type)),
		logging.String("profile", profileName),
		logging.String("correlation_id", cid),
	)

	cred, err := s.factory.Create(credOpts)
	if err != nil {
		s.logger.Error("failed to create credential",
			logging.String("event", "create_credential"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return azcore.AccessToken{}, fmt.Errorf("failed to create credential: %w", err)
	}

	s.logger.Debug("attempting silent token acquisition",
		logging.String("event", "silent_acquire"),
		logging.String("profile", profileName),
		logging.String("correlation_id", cid),
	)

	token, err := cred.GetToken(ctx, options)
	if err != nil {
		if record == (azidentity.AuthenticationRecord{}) {
			s.logger.Warn("no cached record; user not logged in",
				logging.String("event", "silent_acquire"),
				logging.String("correlation_id", cid),
			)
			return azcore.AccessToken{}, errors.New("not logged in")
		}

		s.logger.Warn("silent token acquisition failed",
			logging.String("event", "silent_acquire"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return azcore.AccessToken{}, fmt.Errorf("failed to acquire token silently: %w", err)
	}

	s.logger.Info("silent token acquisition successful",
		logging.String("event", "get_token_success"),
		logging.String("correlation_id", cid),
	)

	return token, nil
}

// Login performs interactive authentication if needed, then retrieves a token.
func (s *AuthService) Login(ctx context.Context, profileName string, opts authdomain.LoginOptions) (*authdomain.LoginResult, error) {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Info("starting login flow",
		logging.String("event", "login_start"),
		logging.String("profile", profileName),
		logging.Bool("force", opts.Force),
		logging.String("correlation_id", cid),
	)

	var (
		record azidentity.AuthenticationRecord
		err    error
	)

	if s.cache != nil {
		s.logger.Debug("loading cached authentication record",
			logging.String("event", "load_cached_record"),
			logging.String("profile", profileName),
			logging.String("correlation_id", cid),
		)

		record, err = s.cache.GetProfile(ctx, profileName)
		if err != nil {
			s.logger.Warn("failed to load cached record",
				logging.String("event", "load_cached_record"),
				logging.Error(err),
				logging.String("correlation_id", cid),
			)
			return nil, fmt.Errorf("failed to retrieve cached auth record: %w", err)
		}
	}

	s.logger.Debug("loading profile configuration",
		logging.String("event", "load_config"),
		logging.String("profile", profileName),
		logging.String("correlation_id", cid),
	)

	cfg, err := s.config.GetConfiguration(ctx, profileName)
	if err != nil {
		s.logger.Error("failed to load configuration",
			logging.String("event", "load_config"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return nil, fmt.Errorf("failed to retrieve configuration: %w", err)
	}

	credOpts := authinfra.CredentialOptions{
		Type:                 cfg.Auth.Type,
		ClientID:             cfg.Auth.ClientID,
		TenantID:             cfg.Auth.TenantID,
		AuthenticationRecord: record,
	}

	s.logger.Debug("creating credential",
		logging.String("event", "create_credential"),
		logging.String("auth_type", string(cfg.Auth.Type)),
		logging.String("correlation_id", cid),
	)

	cred, err := s.factory.Create(credOpts)
	if err != nil {
		s.logger.Error("failed to create credential",
			logging.String("event", "create_credential"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return nil, fmt.Errorf("failed to retrieve credential: %w", err)
	}

	authenticator, ok := cred.(authdomain.Authenticator)
	if !ok {
		s.logger.Error("credential does not support explicit authentication",
			logging.String("event", "authenticator_check"),
			logging.String("correlation_id", cid),
		)
		return nil, errors.New("credential does not support explicit authentication")
	}

	tokenOpts := &policy.TokenRequestOptions{
		Scopes:    opts.Scopes,
		EnableCAE: opts.EnableCAE,
	}

	var token azcore.AccessToken
	maxAttempts := 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		s.logger.Debug("evaluating authentication requirement",
			logging.String("event", "auth_check"),
			logging.Bool("force", opts.Force),
			logging.Bool("has_record", record != (azidentity.AuthenticationRecord{})),
			logging.Int("attempt", attempt),
			logging.String("correlation_id", cid),
		)

		needsAuth := opts.Force || record == (azidentity.AuthenticationRecord{})

		if needsAuth {
			s.logger.Info("performing interactive authentication",
				logging.String("event", "interactive_auth"),
				logging.Int("attempt", attempt),
				logging.String("correlation_id", cid),
			)

			newRecord, err := authenticator.Authenticate(ctx, tokenOpts)
			if err != nil {
				s.logger.Error("interactive authentication failed",
					logging.String("event", "interactive_auth"),
					logging.Error(err),
					logging.String("correlation_id", cid),
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
				s.logger.Error("failed to reload credential",
					logging.String("event", "reload_credential"),
					logging.Error(err),
					logging.String("correlation_id", cid),
				)
				return nil, fmt.Errorf("failed to reload credential: %w", err)
			}
		}

		s.logger.Debug("attempting token acquisition",
			logging.String("event", "acquire_token"),
			logging.Int("attempt", attempt),
			logging.String("correlation_id", cid),
		)

		token, err = cred.GetToken(ctx, *tokenOpts)
		if err != nil {
			if !isAuthRequired(err) {
				s.logger.Error("token acquisition failed",
					logging.String("event", "acquire_token"),
					logging.Error(err),
					logging.Int("attempt", attempt),
					logging.String("correlation_id", cid),
				)
				return nil, fmt.Errorf("failed to retrieve token: %w", err)
			}

			s.logger.Warn("token acquisition requires re-authentication",
				logging.String("event", "acquire_token"),
				logging.Error(err),
				logging.Int("attempt", attempt),
				logging.String("correlation_id", cid),
			)

			record = azidentity.AuthenticationRecord{}
			continue
		}

		break
	}

	if token == (azcore.AccessToken{}) {
		s.logger.Warn("token is empty",
			logging.String("event", "token_empty"),
			logging.String("correlation_id", cid),
		)
		return nil, errors.New("empty token")
	}

	s.logger.Debug("saving authentication record",
		logging.String("event", "save_record"),
		logging.String("profile", profileName),
		logging.String("correlation_id", cid),
	)

	if err := s.cache.SetProfile(ctx, profileName, record); err != nil {
		s.logger.Warn("failed to save authentication record",
			logging.String("event", "save_record"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
	}

	s.logger.Info("login successful",
		logging.String("event", "login_success"),
		logging.String("profile", profileName),
		logging.String("correlation_id", cid),
	)

	return &authdomain.LoginResult{
		AccessToken: token.Token,
		RecordSaved: (record != azidentity.AuthenticationRecord{}),
	}, nil
}

// Logout removes the cached authentication record for the given profile.
// If force is true, logout proceeds even if no record exists.
func (s *AuthService) Logout(ctx context.Context, profileName string, force bool) error {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Info("starting logout flow",
		logging.String("event", "logout_start"),
		logging.String("profile", profileName),
		logging.Bool("force", force),
		logging.String("correlation_id", cid),
	)

	record, err := s.cache.GetProfile(ctx, profileName)
	if err != nil {
		if force {
			s.logger.Warn("failed to load cached record, but force=true; continuing",
				logging.String("event", "load_cached_record"),
				logging.Error(err),
				logging.String("correlation_id", cid),
			)
		} else {
			s.logger.Error("failed to load cached record",
				logging.String("event", "load_cached_record"),
				logging.Error(err),
				logging.String("correlation_id", cid),
			)
			return fmt.Errorf("failed to load cached authentication record: %w", err)
		}
	}

	if record == (azidentity.AuthenticationRecord{}) && !force {
		s.logger.Info("no authentication record found; nothing to do",
			logging.String("event", "logout_noop"),
			logging.String("profile", profileName),
			logging.String("correlation_id", cid),
		)
		return nil
	}

	s.logger.Debug("deleting authentication record",
		logging.String("event", "delete_record"),
		logging.String("profile", profileName),
		logging.String("correlation_id", cid),
	)

	if err := s.cache.DeleteProfile(ctx, profileName); err != nil {
		s.logger.Error("failed to delete authentication record",
			logging.String("event", "delete_record"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return fmt.Errorf("failed to delete authentication record: %w", err)
	}

	s.logger.Info("logout successful",
		logging.String("event", "logout_success"),
		logging.String("profile", profileName),
		logging.String("correlation_id", cid),
	)

	return nil
}
