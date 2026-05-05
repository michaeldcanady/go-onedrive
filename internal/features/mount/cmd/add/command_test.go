package add

import (
	"bytes"
	"context"
	"testing"

	environment "github.com/michaeldcanady/go-onedrive/internal/core/env"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/config"
	drive "github.com/michaeldcanady/go-onedrive/internal/features/drive/domain"
	editor "github.com/michaeldcanady/go-onedrive/internal/features/editor/domain"
	fsdomain "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
	proto "github.com/michaeldcanady/go-onedrive/internal/features/identity/proto"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	"github.com/michaeldcanady/go-onedrive/internal/features/profile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockContainer struct{ mock.Mock }

func (m *mockContainer) Logger() logger.Service     { return m.Called().Get(0).(logger.Service) }
func (m *mockContainer) Config() config.Service     { return m.Called().Get(0).(config.Service) }
func (m *mockContainer) Mounts() mount.Service      { return m.Called().Get(0).(mount.Service) }
func (m *mockContainer) Identity() identity.Service { return m.Called().Get(0).(identity.Service) }
func (m *mockContainer) Profile() profile.Service   { return m.Called().Get(0).(profile.Service) }
func (m *mockContainer) FS() fsdomain.Service       { return m.Called().Get(0).(fsdomain.Service) }
func (m *mockContainer) Environment() environment.Service {
	return m.Called().Get(0).(environment.Service)
}
func (m *mockContainer) Editor() editor.Service { return m.Called().Get(0).(editor.Service) }
func (m *mockContainer) Drive() drive.Service   { return m.Called().Get(0).(drive.Service) }
func (m *mockContainer) URIFactory() *fsdomain.URIFactory {
	return m.Called().Get(0).(*fsdomain.URIFactory)
}

type mockMountService struct{ mock.Mock }

func (m *mockMountService) AddMount(ctx context.Context, cfg mount.MountConfig) error {
	return m.Called(ctx, cfg).Error(0)
}
func (m *mockMountService) RemoveMount(ctx context.Context, path string) error {
	return m.Called(ctx, path).Error(0)
}
func (m *mockMountService) ListMounts(ctx context.Context) ([]mount.MountConfig, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]mount.MountConfig), args.Error(1)
}
func (m *mockMountService) GetMountOptions() map[string][]mount.MountOption {
	return m.Called().Get(0).(map[string][]mount.MountOption)
}
func (m *mockMountService) RegisterValidator(name string, validator mount.OptionValidator) {
	m.Called(name, validator)
}
func (m *mockMountService) RegisterCompletionProvider(name string, provider mount.CompletionProvider) {
	m.Called(name, provider)
}
func (m *mockMountService) GetCompletionProvider(name string) (mount.CompletionProvider, bool) {
	args := m.Called(name)
	return args.Get(0).(mount.CompletionProvider), args.Bool(1)
}

type mockIdentityService struct{ mock.Mock }

func (m *mockIdentityService) RegisterAuthenticator(p string, a identity.Authenticator) {
	m.Called(p, a)
}
func (m *mockIdentityService) RegisterAuthorizer(p string, a identity.Authorizer) { m.Called(p, a) }
func (m *mockIdentityService) GetAuthenticator(p string) (identity.Authenticator, error) {
	args := m.Called(p)
	return args.Get(0).(identity.Authenticator), args.Error(1)
}
func (m *mockIdentityService) Authenticate(ctx context.Context, p string, req *proto.AuthenticateRequest) (*proto.AuthenticateResponse, error) {
	args := m.Called(ctx, p, req)
	return args.Get(0).(*proto.AuthenticateResponse), args.Error(1)
}
func (m *mockIdentityService) Login(ctx context.Context, p string, opts identity.LoginOptions) (*proto.AuthenticateResponse, error) {
	args := m.Called(ctx, p, opts)
	return args.Get(0).(*proto.AuthenticateResponse), args.Error(1)
}
func (m *mockIdentityService) Logout(ctx context.Context, p string, id string) error {
	return m.Called(ctx, p, id).Error(0)
}
func (m *mockIdentityService) Token(ctx context.Context, p string, req *proto.GetTokenRequest) (*proto.GetTokenResponse, error) {
	args := m.Called(ctx, p, req)
	return args.Get(0).(*proto.GetTokenResponse), args.Error(1)
}
func (m *mockIdentityService) GetStore() identity.AccountStore {
	return m.Called().Get(0).(identity.AccountStore)
}
func (m *mockIdentityService) GetAccount(ctx context.Context, id string) (*identity.Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*identity.Account), args.Error(1)
}
func (m *mockIdentityService) ListProviders() []string { return m.Called().Get(0).([]string) }

type mockLoggerService struct{ mock.Mock }

func (m *mockLoggerService) CreateLogger(name string) (logger.Logger, error) {
	args := m.Called(name)
	return args.Get(0).(logger.Logger), args.Error(1)
}
func (m *mockLoggerService) SetAllLevel(level logger.Level) { m.Called(level) }
func (m *mockLoggerService) Reconfigure(level logger.Level, output string, format string) error {
	return m.Called(level, output, format).Error(0)
}

type mockLogger struct{ mock.Mock }

func (m *mockLogger) Debug(msg string, fields ...logger.Field) { m.Called(msg, fields) }
func (m *mockLogger) Info(msg string, fields ...logger.Field)  { m.Called(msg, fields) }
func (m *mockLogger) Warn(msg string, fields ...logger.Field)  { m.Called(msg, fields) }
func (m *mockLogger) Error(msg string, fields ...logger.Field) { m.Called(msg, fields) }
func (m *mockLogger) SetLevel(level logger.Level)              { m.Called(level) }
func (m *mockLogger) With(fields ...logger.Field) logger.Logger {
	args := m.Called(fields)
	return args.Get(0).(logger.Logger)
}
func (m *mockLogger) WithContext(ctx context.Context) logger.Logger {
	args := m.Called(ctx)
	return args.Get(0).(logger.Logger)
}

type mockVFS struct {
	mock.Mock
	fsdomain.Service
}

func (m *mockVFS) Resolve(absPath string) (string, string, error) {
	args := m.Called(absPath)
	return args.String(0), args.String(1), args.Error(2)
}

func TestAddCmd_Integration(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		setup   func(m *mockContainer, mMount *mockMountService, mIdent *mockIdentityService, mLogSvc *mockLoggerService, mLog *mockLogger)
		wantErr bool
	}{
		{
			name: "add mount success",
			args: []string{"/od", "onedrive", "user1"},
			setup: func(m *mockContainer, mMount *mockMountService, mIdent *mockIdentityService, mLogSvc *mockLoggerService, mLog *mockLogger) {
				mLogSvc.On("CreateLogger", "mount-add").Return(mLog, nil)
				m.On("Logger").Return(mLogSvc)
				m.On("Mounts").Return(mMount)
				m.On("Identity").Return(mIdent)

				mVFS := new(mockVFS)
				mVFS.On("Resolve", "/od").Return("/od", "", nil)
				m.On("URIFactory").Return(fsdomain.NewURIFactory(mVFS))

				mMount.On("GetMountOptions").Return(map[string][]mount.MountOption{})
				mIdent.On("GetAccount", mock.Anything, "user1").Return(&identity.Account{ID: "user1"}, nil)
				mMount.On("AddMount", mock.Anything, mock.MatchedBy(func(cfg mount.MountConfig) bool {
					return cfg.Path == "/od:" && cfg.Type == "onedrive" && cfg.IdentityID == "user1"
				})).Return(nil)

				mLog.On("WithContext", mock.Anything).Return(mLog)
				mLog.On("Info", mock.Anything, mock.Anything).Return()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mContainer := new(mockContainer)
			mMount := new(mockMountService)
			mIdent := new(mockIdentityService)
			mLogSvc := new(mockLoggerService)
			mLog := new(mockLogger)

			tt.setup(mContainer, mMount, mIdent, mLogSvc, mLog)

			cmd := CreateAddCmd(mContainer)
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
