package list

import (
	"bytes"
	"testing"

	environment "github.com/michaeldcanady/go-onedrive/internal/core/env"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/config"
	drive "github.com/michaeldcanady/go-onedrive/internal/features/drive/domain"
	editor "github.com/michaeldcanady/go-onedrive/internal/features/editor/domain"
	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
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
func (m *mockContainer) FS() fs.Service             { return m.Called().Get(0).(fs.Service) }
func (m *mockContainer) Environment() environment.Service {
	return m.Called().Get(0).(environment.Service)
}
func (m *mockContainer) Editor() editor.Service     { return m.Called().Get(0).(editor.Service) }
func (m *mockContainer) Drive() drive.Service       { return m.Called().Get(0).(drive.Service) }
func (m *mockContainer) URIFactory() *fs.URIFactory { return m.Called().Get(0).(*fs.URIFactory) }

type mockLoggerService struct{ mock.Mock }

func (m *mockLoggerService) CreateLogger(name string) (logger.Logger, error) {
	args := m.Called(name)
	return args.Get(0).(logger.Logger), args.Error(1)
}
func (m *mockLoggerService) SetAllLevel(level logger.Level) { m.Called(level) }
func (m *mockLoggerService) Reconfigure(level logger.Level, output string, format string) error {
	return m.Called(level, output, format).Error(0)
}

func TestListCommand_Integration(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		setup   func(m *mockContainer, mSvc *mockDriveService, mLogSvc *mockLoggerService, mLog *mockLogger)
		wantErr bool
	}{
		{
			name: "list drives success",
			args: []string{"list"},
			setup: func(m *mockContainer, mSvc *mockDriveService, mLogSvc *mockLoggerService, mLog *mockLogger) {
				mLogSvc.On("CreateLogger", "drive-list").Return(mLog, nil)
				m.On("Logger").Return(mLogSvc)
				m.On("Drive").Return(mSvc)

				mSvc.On("ListDrives", mock.Anything, "").Return([]drive.Drive{
					{ID: "d1", Name: "Drive 1", Type: "personal"},
				}, nil)

				mLog.On("WithContext", mock.Anything).Return(mLog)
				mLog.On("Debug", mock.Anything, mock.Anything).Return()
				mLog.On("Info", mock.Anything, mock.Anything).Return()
			},
			wantErr: false,
		},
		{
			name: "list drives with identity",
			args: []string{"list", "--id", "user@example.com"},
			setup: func(m *mockContainer, mSvc *mockDriveService, mLogSvc *mockLoggerService, mLog *mockLogger) {
				mLogSvc.On("CreateLogger", "drive-list").Return(mLog, nil)
				m.On("Logger").Return(mLogSvc)
				m.On("Drive").Return(mSvc)

				mSvc.On("ListDrives", mock.Anything, "user@example.com").Return([]drive.Drive{
					{ID: "d1", Name: "Drive 1", Type: "personal"},
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
			mContainer := new(mockContainer)
			mSvc := new(mockDriveService)
			mLogSvc := new(mockLoggerService)
			mLog := new(mockLogger)

			tt.setup(mContainer, mSvc, mLogSvc, mLog)

			cmd := CreateListCmd(mContainer)
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args[1:])

			err := cmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
