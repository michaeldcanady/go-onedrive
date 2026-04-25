package identity

import (
	"context"
	"testing"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	proto "github.com/michaeldcanady/go-onedrive/internal/features/identity/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthenticator struct {
	mock.Mock
}

func (m *mockAuthenticator) ProviderName() string { return "test-provider" }
func (m *mockAuthenticator) Logout(ctx context.Context, identityID string) error {
	return m.Called(ctx, identityID).Error(0)
}
func (m *mockAuthenticator) Authenticate(ctx context.Context, req *proto.AuthenticateRequest) (*proto.AuthenticateResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*proto.AuthenticateResponse), args.Error(1)
}

type mockAccountStore struct {
	mock.Mock
}

func (m *mockAccountStore) Save(ctx context.Context, provider string, token AccessToken) error {
	return m.Called(ctx, provider, token).Error(0)
}
func (m *mockAccountStore) Get(ctx context.Context, provider string, accountID string) (AccessToken, error) {
	args := m.Called(ctx, provider, accountID)
	return args.Get(0).(AccessToken), args.Error(1)
}
func (m *mockAccountStore) List(ctx context.Context, provider string) ([]string, error) {
	args := m.Called(ctx, provider)
	return args.Get(0).([]string), args.Error(1)
}
func (m *mockAccountStore) Delete(ctx context.Context, provider string, accountID string) error {
	return m.Called(ctx, provider, accountID).Error(0)
}
func (m *mockAccountStore) Close() error { return m.Called().Error(0) }

func TestIdentityService_Functional(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(mAuth *mockAuthenticator, mStore *mockAccountStore)
		action  func(svc Service) error
		wantErr bool
	}{
		{
			name: "complete login and token flow",
			setup: func(mAuth *mockAuthenticator, mStore *mockAccountStore) {
				// Mock successful authentication
				resp := &proto.AuthenticateResponse{
					Identity: &proto.Identity{Id: "user1", Email: "user1@example.com"},
					Token:    &proto.AccessToken{Token: "token1", ExpiresAt: time.Now().Add(1 * time.Hour).Unix()},
				}
				mAuth.On("Authenticate", mock.Anything, mock.Anything).Return(resp, nil)

				// Mock successful persistence
				mStore.On("Save", mock.Anything, "test-provider", mock.MatchedBy(func(token AccessToken) bool {
					return token.Token == "token1" && token.AccountID == "user1"
				})).Return(nil)

				// Mock token retrieval (cache hit)
				mStore.On("Get", mock.Anything, "test-provider", "user1").Return(AccessToken{
					Token:     "token1",
					AccountID: "user1",
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}, nil)
			},
			action: func(svc Service) error {
				// 1. Perform Login
				_, err := svc.Login(ctx, "test-provider", LoginOptions{AccountID: "user1"})
				if err != nil {
					return err
				}

				// 2. Retrieve Token (should hit cache)
				_, err = svc.Token(ctx, "test-provider", &proto.GetTokenRequest{IdentityId: "user1"})
				return err
			},
			wantErr: false,
		},
		{
			name: "logout flow",
			setup: func(mAuth *mockAuthenticator, mStore *mockAccountStore) {
				mAuth.On("Logout", mock.Anything, "user1").Return(nil)
			},
			action: func(svc Service) error {
				return svc.Logout(ctx, "test-provider", "user1")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mAuth := new(mockAuthenticator)
			mStore := new(mockAccountStore)
			tt.setup(mAuth, mStore)

			registry := NewRegistry(mStore, &dummyLogger{})
			registry.RegisterAuthenticator("test-provider", mAuth)

			err := tt.action(registry)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mAuth.AssertExpectations(t)
			mStore.AssertExpectations(t)
		})
	}
}

type dummyLogger struct{}

func (l *dummyLogger) Debug(msg string, fields ...logger.Field) {}
func (l *dummyLogger) Error(msg string, fields ...logger.Field) {}
