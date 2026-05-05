package list

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

func TestListCmd_Integration(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		setup   func(m *mockContainer, mMount *mockMountService, mLogSvc *mockLoggerService, mLog *mockLogger)
		wantErr bool
	}{
		{
			name: "list mounts success",
			args: []string{"-f", "json"},
			setup: func(m *mockContainer, mMount *mockMountService, mLogSvc *mockLoggerService, mLog *mockLogger) {
				mLogSvc.On("CreateLogger", "mount-list").Return(mLog, nil)
				m.On("Logger").Return(mLogSvc)
				m.On("Mounts").Return(mMount)

				mMount.On("ListMounts", mock.Anything).Return([]mount.MountConfig{
					{Path: "/od", Type: "onedrive", IdentityID: "user1"},
				}, nil)

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
			mLogSvc := new(mockLoggerService)
			mLog := new(mockLogger)

			tt.setup(mContainer, mMount, mLogSvc, mLog)

			cmd := CreateListCmd(mContainer)
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
