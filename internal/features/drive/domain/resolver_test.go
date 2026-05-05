package drive

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDriveService struct {
	mock.Mock
}

func (m *mockDriveService) ListDrives(ctx context.Context, identityID string) ([]Drive, error) {
	args := m.Called(ctx, identityID)
	return args.Get(0).([]Drive), args.Error(1)
}

func (m *mockDriveService) ResolveDrive(ctx context.Context, driveRef string, identityID string) (Drive, error) {
	args := m.Called(ctx, driveRef, identityID)
	return args.Get(0).(Drive), args.Error(1)
}

func (m *mockDriveService) ResolvePersonalDrive(ctx context.Context, identityID string) (Drive, error) {
	args := m.Called(ctx, identityID)
	return args.Get(0).(Drive), args.Error(1)
}

func TestDefaultResolver_GetActiveDriveID(t *testing.T) {
	tests := []struct {
		name       string
		identityID string
		setup      func(m *mockDriveService)
		want       string
		wantErr    bool
	}{
		{
			name:       "success",
			identityID: "user1",
			setup: func(m *mockDriveService) {
				m.On("ResolvePersonalDrive", mock.Anything, "user1").Return(Drive{ID: "d1"}, nil)
			},
			want:    "d1",
			wantErr: false,
		},
		{
			name:       "failure",
			identityID: "user1",
			setup: func(m *mockDriveService) {
				m.On("ResolvePersonalDrive", mock.Anything, "user1").Return(Drive{}, assert.AnError)
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mSvc := new(mockDriveService)
			tt.setup(mSvc)

			resolver := NewDefaultResolver(mSvc, tt.identityID)
			got, err := resolver.GetActiveDriveID(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			mSvc.AssertExpectations(t)
		})
	}
}
