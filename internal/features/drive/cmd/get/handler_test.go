package get

import (
	"bytes"
	"context"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/drive/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDriveService struct {
	mock.Mock
}

func (m *mockDriveService) ListDrives(ctx context.Context, identityID string) ([]drive.Drive, error) {
	args := m.Called(ctx, identityID)
	return args.Get(0).([]drive.Drive), args.Error(1)
}

func (m *mockDriveService) ResolveDrive(ctx context.Context, driveRef string, identityID string) (drive.Drive, error) {
	args := m.Called(ctx, driveRef, identityID)
	return args.Get(0).(drive.Drive), args.Error(1)
}

func (m *mockDriveService) ResolvePersonalDrive(ctx context.Context, identityID string) (drive.Drive, error) {
	args := m.Called(ctx, identityID)
	return args.Get(0).(drive.Drive), args.Error(1)
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
		name     string
		driveRef string
		setup    func(m *mockDriveService, l *mockLogger)
		wantErr  bool
	}{
		{
			name:     "get drive success",
			driveRef: "personal",
			setup: func(m *mockDriveService, l *mockLogger) {
				m.On("ResolveDrive", mock.Anything, "personal", "").Return(drive.Drive{
					ID: "d1", Name: "Personal Drive", Type: "personal",
				}, nil)
				l.On("WithContext", mock.Anything).Return(l)
				l.On("Debug", mock.Anything, mock.Anything).Return()
				l.On("Info", mock.Anything, mock.Anything).Return()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mSvc := new(mockDriveService)
			mLog := new(mockLogger)
			tt.setup(mSvc, mLog)

			handler := NewCommand(mSvc, mLog)

			ctx := &CommandContext{
				Ctx: context.Background(),
				Options: Options{
					DriveRef: tt.driveRef,
					Stdout:   new(bytes.Buffer),
				},
			}

			err := handler.Execute(ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mSvc.AssertExpectations(t)
		})
	}
}
