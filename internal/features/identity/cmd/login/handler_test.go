package login

import (
	"bytes"
	"context"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/config"
	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
	proto "github.com/michaeldcanady/go-onedrive/internal/features/identity/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockConfigService struct {
	mock.Mock
}

func (m *mockConfigService) GetConfig(ctx context.Context) (config.Config, error) {
	args := m.Called(ctx)
	return args.Get(0).(config.Config), args.Error(1)
}

func (m *mockConfigService) GetPath(ctx context.Context) (string, bool) {
	args := m.Called(ctx)
	return args.String(0), args.Bool(1)
}

func (m *mockConfigService) SaveConfig(ctx context.Context, cfg config.Config) error {
	return m.Called(ctx, cfg).Error(0)
}

func (m *mockConfigService) UpdateConfig(ctx context.Context, key string, value string) error {
	return m.Called(ctx, key, value).Error(0)
}

func (m *mockConfigService) SetOverride(ctx context.Context, path string) error {
	return m.Called(ctx, path).Error(0)
}

type mockIdentityService struct {
	mock.Mock
}

func (m *mockIdentityService) RegisterAuthenticator(provider string, auth identity.Authenticator) {
	m.Called(provider, auth)
}
func (m *mockIdentityService) RegisterAuthorizer(provider string, auth identity.Authorizer) {
	m.Called(provider, auth)
}
func (m *mockIdentityService) GetAuthenticator(provider string) (identity.Authenticator, error) {
	args := m.Called(provider)
	return args.Get(0).(identity.Authenticator), args.Error(1)
}
func (m *mockIdentityService) Authenticate(ctx context.Context, provider string, req *proto.AuthenticateRequest) (*proto.AuthenticateResponse, error) {
	args := m.Called(ctx, provider, req)
	return args.Get(0).(*proto.AuthenticateResponse), args.Error(1)
}
func (m *mockIdentityService) Login(ctx context.Context, provider string, opts identity.LoginOptions) (*proto.AuthenticateResponse, error) {
	args := m.Called(ctx, provider, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*proto.AuthenticateResponse), args.Error(1)
}
func (m *mockIdentityService) Logout(ctx context.Context, provider string, identityID string) error {
	return m.Called(ctx, provider, identityID).Error(0)
}
func (m *mockIdentityService) Token(ctx context.Context, provider string, req *proto.GetTokenRequest) (*proto.GetTokenResponse, error) {
	args := m.Called(ctx, provider, req)
	return args.Get(0).(*proto.GetTokenResponse), args.Error(1)
}
func (m *mockIdentityService) GetStore() identity.AccountStore {
	return m.Called().Get(0).(identity.AccountStore)
}
func (m *mockIdentityService) GetAccount(ctx context.Context, identityID string) (*identity.Account, error) {
	args := m.Called(ctx, identityID)
	return args.Get(0).(*identity.Account), args.Error(1)
}
func (m *mockIdentityService) ListProviders() []string {
	return m.Called().Get(0).([]string)
}

type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) Debug(msg string, fields ...logger.Field) { m.Called(msg, fields) }
func (m *mockLogger) Info(msg string, fields ...logger.Field)  { m.Called(msg, fields) }
func (m *mockLogger) Warn(msg string, fields ...logger.Field)  { m.Called(msg, fields) }
func (m *mockLogger) Error(msg string, fields ...logger.Field) { m.Called(msg, fields) }
func (m *mockLogger) SetLevel(level logger.Level)              { m.Called(level) }
func (m *mockLogger) WithContext(ctx context.Context) logger.Logger {
	return m.Called(ctx).Get(0).(logger.Logger)
}
func (m *mockLogger) With(fields ...logger.Field) logger.Logger {
	return m.Called(fields).Get(0).(logger.Logger)
}

func TestHandler_Execute(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(mCfg *mockConfigService, mIdent *mockIdentityService, mLog *mockLogger)
		wantErr bool
	}{
		{
			name: "login success",
			setup: func(mCfg *mockConfigService, mIdent *mockIdentityService, mLog *mockLogger) {
				mCfg.On("GetConfig", mock.Anything).Return(config.Config{
					Auth: config.AuthenticationConfig{Provider: "microsoft"},
				}, nil)
				mIdent.On("Login", mock.Anything, "microsoft", mock.Anything).Return(&proto.AuthenticateResponse{
					Identity: &proto.Identity{Id: "user1", Email: "user1@example.com"},
				}, nil)
				mLog.On("WithContext", mock.Anything).Return(mLog)
				mLog.On("Debug", mock.Anything, mock.Anything).Return()
				mLog.On("Info", mock.Anything, mock.Anything).Return()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mCfg := new(mockConfigService)
			mIdent := new(mockIdentityService)
			mLog := new(mockLogger)
			tt.setup(mCfg, mIdent, mLog)

			handler := NewCommand(mCfg, mIdent, mLog)

			ctx := &CommandContext{
				Ctx: context.Background(),
				Options: &Options{
					Stdout: new(bytes.Buffer),
				},
			}

			err := handler.Execute(ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mIdent.AssertExpectations(t)
		})
	}
}
