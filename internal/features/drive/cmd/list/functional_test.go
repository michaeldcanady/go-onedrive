package list

import (
	"bytes"
	"context"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/features/drive/domain"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	formatting "github.com/michaeldcanady/go-onedrive/pkg/format"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDriveSource struct {
	mock.Mock
}

func (m *mockDriveSource) ListDrives(ctx context.Context, provider string) ([]fs.Drive, error) {
	args := m.Called(ctx, provider)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]fs.Drive), args.Error(1)
}

func (m *mockDriveSource) GetPersonalDrive(ctx context.Context, provider string) (fs.Drive, error) {
	args := m.Called(ctx, provider)
	return args.Get(0).(fs.Drive), args.Error(1)
}

type mockMountProvider struct {
	mock.Mock
}

func (m *mockMountProvider) ListMounts(ctx context.Context) ([]mount.MountConfig, error) {
	args := m.Called(ctx)
	return args.Get(0).([]mount.MountConfig), args.Error(1)
}

func TestList_Functional(t *testing.T) {
	tests := []struct {
		name     string
		identity string
		setup    func(mSource *mockDriveSource, mMounts *mockMountProvider, mLog *mockLogger)
		wantErr  bool
		expected string
	}{
		{
			name: "list all drives successfully",
			setup: func(mSource *mockDriveSource, mMounts *mockMountProvider, mLog *mockLogger) {
				mMounts.On("ListMounts", mock.Anything).Return([]mount.MountConfig{
					{Path: "/od", IdentityID: "user1"},
				}, nil)
				mSource.On("ListDrives", mock.Anything, "/od").Return([]fs.Drive{
					{ID: "d1", Name: "My Drive", Type: "personal"},
				}, nil)
				mLog.On("WithContext", mock.Anything).Return(mLog)
				mLog.On("Debug", mock.Anything, mock.Anything).Return()
				mLog.On("Info", mock.Anything, mock.Anything).Return()
			},
			wantErr:  false,
			expected: "My Drive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mSource := new(mockDriveSource)
			mMounts := new(mockMountProvider)
			mLog := new(mockLogger)
			tt.setup(mSource, mMounts, mLog)

			svc := drive.NewDefaultService(mSource, mMounts, mLog)
			formatterFactory := formatting.NewFormatterFactory()
			handler := NewCommand(svc, formatterFactory, mLog)

			buf := new(bytes.Buffer)
			cmdCtx := &CommandContext{
				Ctx:    ctx,
				Format: formatting.FormatTable,
				Options: Options{
					IdentityID: tt.identity,
					Format:     "table",
					Stdout:     buf,
				},
			}

			err := handler.Validate(cmdCtx)
			assert.NoError(t, err)

			err = handler.Execute(cmdCtx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, buf.String(), tt.expected)
			}
		})
	}
}
