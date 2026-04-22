package identity

import (
	"context"
	"fmt"
	"sync"
	"time"

	proto "github.com/michaeldcanady/go-onedrive/internal/features/identity/proto"
	"github.com/michaeldcanady/go-onedrive/internal/features/logger"
)

// Logger defines the interface required for logging within the identity service.
type Logger interface {
	Debug(msg string, fields ...logger.Field)
	Error(msg string, fields ...logger.Field)
}

// Service defines the interface for managing identity providers.
type Service interface {
	RegisterAuthenticator(provider string, auth Authenticator)
	RegisterAuthorizer(provider string, auth Authorizer)
	GetAuthenticator(provider string) (Authenticator, error)
	Authenticate(ctx context.Context, provider string, req *proto.AuthenticateRequest) (*proto.AuthenticateResponse, error)
	Token(ctx context.Context, provider string, req *proto.GetTokenRequest) (*proto.GetTokenResponse, error)
	GetStore() AccountStore
}

func (r *Registry) GetAuthenticator(provider string) (Authenticator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	auth, ok := r.authenticators[provider]
	if !ok {
		return nil, fmt.Errorf("authenticator for provider %s not found", provider)
	}
	return auth, nil
}

// Registry is a thread-safe implementation of the Service interface, acting as the mediator.
type Registry struct {
	mu             sync.RWMutex
	authenticators map[string]Authenticator
	authorizers    map[string]Authorizer
	store          AccountStore
	logger         Logger
}

// NewRegistry initializes a new instance of the Registry with persistence and logging.
func NewRegistry(store AccountStore, logger Logger) *Registry {
	return &Registry{
		authenticators: make(map[string]Authenticator),
		authorizers:    make(map[string]Authorizer),
		store:          store,
		logger:         logger,
	}
}

func (r *Registry) RegisterAuthenticator(provider string, auth Authenticator) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.authenticators[provider] = auth
}

func (r *Registry) RegisterAuthorizer(provider string, auth Authorizer) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.authorizers[provider] = auth
}

func (r *Registry) Authenticate(ctx context.Context, provider string, req *proto.AuthenticateRequest) (*proto.AuthenticateResponse, error) {
	r.mu.RLock()
	auth, ok := r.authenticators[provider]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("authenticator for provider %s not found", provider)
	}

	resp, err := auth.Authenticate(ctx, req)
	if err != nil {
		r.logger.Error("Authentication failed", logger.String("provider", provider), logger.Error(err))
		return nil, err
	}

	return resp, nil
}

func (r *Registry) Token(ctx context.Context, provider string, req *proto.GetTokenRequest) (*proto.GetTokenResponse, error) {
	accountID := req.GetIdentityId()
	if accountID == "" {
		ids, err := r.store.List(ctx, provider)
		if err == nil && len(ids) > 0 {
			accountID = ids[0]
		}
	}

	// 1. Check cache
	if token, err := r.store.Get(ctx, provider, accountID); err == nil {
		if token.ExpiresAt.After(time.Now().Add(5 * time.Minute)) {
			r.logger.Debug("using valid cached token from repository", logger.String("provider", provider), logger.String("identity", accountID))
			return &proto.GetTokenResponse{Token: ToProtoAccessToken(token)}, nil
		}
	}

	// 2. Delegate to authorizer (plugin)
	r.mu.RLock()
	authorizer, ok := r.authorizers[provider]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("authorizer for provider %s not found", provider)
	}

	// Create a new request with the resolved accountID
	newReq := &proto.GetTokenRequest{
		IdentityId: accountID,
		Scopes:     req.GetScopes(),
	}

	resp, err := authorizer.Token(ctx, newReq)
	if err != nil {
		r.logger.Error("Token retrieval failed", logger.String("provider", provider), logger.Error(err))
		return nil, err
	}

	// 3. Persist new token
	token := FromProtoAccessToken(resp.GetToken(), accountID)
	token.AccountID = accountID
	if err := r.store.Save(ctx, provider, token); err != nil {
		r.logger.Error("Failed to save new access token", logger.Error(err))
	}

	return resp, nil
}

func (r *Registry) GetStore() AccountStore {
	return r.store
}
