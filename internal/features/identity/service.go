package identity

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/core/plugins"
	identity_proto "github.com/michaeldcanady/go-onedrive/internal/features/plugins/proto/identity"
	"github.com/pkg/browser"
)

type IdentityService struct {
	repo          Repository
	pluginManager plugins.Manager
	tokenService  TokenService
	logger        logger.Service
}

// NewIdentityService returns a new [*IdentityService] initialized with required dependencies.
// It leverages [plugins.Manager] to communicate with provider-specific login logic.
func NewIdentityService(repo Repository, pm plugins.Manager, ts TokenService, l logger.Service) *IdentityService {
	return &IdentityService{
		repo:          repo,
		pluginManager: pm,
		tokenService:  ts,
		logger:        l,
	}
}

func (s *IdentityService) Login(ctx context.Context, provider string, options map[string]string) (*Identity, error) {
	l := logger.WithContext(s.logger, ctx)

	pluginName := fmt.Sprintf("identity-%s", provider)
	client, err := s.pluginManager.GetIdentityPlugin(pluginName)
	if err != nil {
		return nil, fmt.Errorf("failed to get identity plugin %s: %w", pluginName, err)
	}

	stream, err := client.Login(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open login stream: %w", err)
	}

	// 1. Send configuration
	err = stream.Send(&identity_proto.LoginRequest{
		Payload: &identity_proto.LoginRequest_Config{
			Config: &identity_proto.Config{
				Options: options,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send config: %w", err)
	}

	// 2. Interaction Loop
	for {
		resp, err := stream.Recv()
		if err != nil {
			return nil, fmt.Errorf("stream receive failed: %w", err)
		}

		// Handle Interaction Request
		if req := resp.GetInteractionRequest(); req != nil {
			if msg := req.GetDisplayMessage(); msg != nil {
				fmt.Println(msg.Message)
			} else if open := req.GetOpenUrl(); open != nil {
				fmt.Printf("Please open your browser at: %s\n", open.Url)
				if err := browser.OpenURL(open.Url); err != nil {
					l.Warn("failed to open browser", "url", open.Url, "error", err)
				}
			}

			// Send acknowledgement back to plugin
			err = stream.Send(&identity_proto.LoginRequest{
				Payload: &identity_proto.LoginRequest_InteractionResponse{
					InteractionResponse: &identity_proto.InteractionResponse{},
				},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to send interaction response: %w", err)
			}
			continue
		}

		// Handle Final Result
		if result := resp.GetResult(); result != nil {
			identity := FromProtoIdentity(result.Identity)
			token := FromProtoToken(result.Token)

			if err := s.repo.SaveIdentity(identity); err != nil {
				return nil, fmt.Errorf("failed to save identity: %w", err)
			}

			if err := s.tokenService.SaveToken(ctx, provider, identity.ID, token); err != nil {
				return nil, fmt.Errorf("failed to save token: %w", err)
			}

			l.Info("identity logged in", "provider", provider, "identity", identity.ID)
			return identity, nil
		}

		return nil, fmt.Errorf("unexpected message from plugin")
	}
}

func (s *IdentityService) Logout(ctx context.Context, identityID string) error {
	l := logger.WithContext(s.logger, ctx)

	i, err := s.repo.GetIdentity(identityID)
	if err != nil {
		return err
	}
	if i == nil {
		return fmt.Errorf("identity not found: %s", identityID)
	}

	// Optional: call plugin logout
	pluginName := fmt.Sprintf("identity-%s", i.Provider)
	client, err := s.pluginManager.GetIdentityPlugin(pluginName)
	if err == nil {
		if _, err := client.Logout(ctx, &identity_proto.LogoutRequest{
			IdentityId: identityID,
		}); err != nil {
			l.Warn("failed to logout from plugin", "plugin", pluginName, "identity", identityID, "error", err)
		}
	}

	if err := s.tokenService.SaveToken(ctx, i.Provider, identityID, &Token{}); err != nil {
		return err
	}

	l.Info("identity logged out", "identity", identityID)
	return nil
}

func (s *IdentityService) List(ctx context.Context) ([]*Identity, error) {
	return s.repo.ListIdentities()
}

func (s *IdentityService) GetIdentity(ctx context.Context, identityID string) (*Identity, error) {
	return s.repo.GetIdentity(identityID)
}

func (s *IdentityService) FindIdentity(ctx context.Context, query string) (*Identity, error) {
	identities, err := s.repo.ListIdentities()
	if err != nil {
		return nil, err
	}

	for _, iden := range identities {
		if iden.ID == query || iden.Email == query || iden.DisplayName == query {
			return iden, nil
		}
	}

	return nil, fmt.Errorf("identity not found: %s", query)
}
