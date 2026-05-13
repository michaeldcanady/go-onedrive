package identity

import (
	"context"
	"fmt"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/core/plugins"
	"github.com/michaeldcanady/go-onedrive/internal/features/config"
	identity_proto "github.com/michaeldcanady/go-onedrive/internal/features/plugins/proto/identity"
)

type DefaultTokenService struct {
	repo          Repository
	pluginManager plugins.Manager
	config        config.Service
	logger        logger.Service
}

// NewTokenService returns a new [*DefaultTokenService] initialized with required dependencies.
// It uses [config.Service] to retrieve client credentials required for token refreshing.
func NewTokenService(repo Repository, pm plugins.Manager, cs config.Service, l logger.Service) *DefaultTokenService {
	return &DefaultTokenService{
		repo:          repo,
		pluginManager: pm,
		config:        cs,
		logger:        l,
	}
}

func (s *DefaultTokenService) GetToken(ctx context.Context, provider string, identityID string) (*Token, error) {
	l := logger.WithContext(s.logger, ctx)

	token, err := s.repo.GetToken(provider, identityID)
	if err != nil {
		return nil, err
	}
	if token == nil || token.AccessToken == "" {
		return nil, fmt.Errorf("token not found for %s:%s", provider, identityID)
	}

	// Proactive refresh if token is expired or near expiration (e.g., 5 mins)
	if time.Until(token.ExpiresAt) < 5*time.Minute {
		l.Info("token near expiration, refreshing", "provider", provider, "identity", identityID)
		return s.RefreshToken(ctx, provider, identityID)
	}

	return token, nil
}

func (s *DefaultTokenService) SaveToken(ctx context.Context, provider string, identityID string, token *Token) error {
	return s.repo.SaveToken(provider, identityID, token)
}

func (s *DefaultTokenService) RefreshToken(ctx context.Context, provider string, identityID string) (*Token, error) {
	token, err := s.repo.GetToken(provider, identityID)
	if err != nil {
		return nil, err
	}
	if token == nil || token.RefreshToken == "" {
		return nil, fmt.Errorf("cannot refresh: no refresh token found for %s:%s. Please run 'odc identity login --provider %s'", provider, identityID, provider)
	}

	options := make(map[string]string)
	clientIDKey := fmt.Sprintf("identity.%s.client_id", provider)
	clientSecretKey := fmt.Sprintf("identity.%s.client_secret", provider)

	if val, err := s.config.Get(clientIDKey); err == nil && val != nil {
		options["client_id"] = fmt.Sprintf("%v", val)
	}
	if val, err := s.config.Get(clientSecretKey); err == nil && val != nil {
		options["client_secret"] = fmt.Sprintf("%v", val)
	}

	pluginName := fmt.Sprintf("identity-%s", provider)
	client, err := s.pluginManager.GetIdentityPlugin(pluginName)
	if err != nil {
		return nil, fmt.Errorf("failed to get identity plugin %s: %w", pluginName, err)
	}

	resp, err := client.Refresh(ctx, &identity_proto.RefreshRequest{
		RefreshToken: token.RefreshToken,
		Options:      options,
	})
	if err != nil {
		return nil, fmt.Errorf("token refresh failed: %w. Your session may have expired, please run 'odc identity login --provider %s' to re-authenticate", err, provider)
	}

	newToken := FromProtoToken(resp.Token)
	// Preserve the old refresh token if the new one is empty
	if newToken.RefreshToken == "" {
		newToken.RefreshToken = token.RefreshToken
	}

	if err := s.repo.SaveToken(provider, identityID, newToken); err != nil {
		return nil, fmt.Errorf("failed to save refreshed token: %w", err)
	}

	return newToken, nil
}
