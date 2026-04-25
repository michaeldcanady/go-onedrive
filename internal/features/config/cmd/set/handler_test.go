package set

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/config"
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
		key     string
		value   string
		setup   func(m *mockConfigService, l *mockLogger)
		wantErr bool
	}{
		{
			name:  "set config success",
			key:   "logging.level",
			value: "debug",
			setup: func(m *mockConfigService, l *mockLogger) {
				m.On("UpdateConfig", mock.Anything, "logging.level", "debug").Return(nil)
				l.On("WithContext", mock.Anything).Return(l)
				l.On("Debug", mock.Anything, mock.Anything).Return()
				l.On("Info", mock.Anything, mock.Anything).Return()
			},
			wantErr: false,
		},
		{
			name:  "config service error",
			key:   "invalid.key",
			value: "val",
			setup: func(m *mockConfigService, l *mockLogger) {
				m.On("UpdateConfig", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("service error"))
				l.On("WithContext", mock.Anything).Return(l)
				l.On("Debug", mock.Anything, mock.Anything).Return()
				l.On("Error", mock.Anything, mock.Anything).Return()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mSvc := new(mockConfigService)
			mLog := new(mockLogger)
			tt.setup(mSvc, mLog)

			handler := NewCommand(mSvc, mLog)

			ctx := &CommandContext{
				Ctx: context.Background(),
				Options: &Options{
					Key:    tt.key,
					Value:  tt.value,
					Stdout: new(bytes.Buffer),
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
