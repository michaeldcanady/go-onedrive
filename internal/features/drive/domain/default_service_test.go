package drive

import (
	"context"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type dummyLogger struct{}

func (l *dummyLogger) Debug(msg string, fields ...logger.Field) {}
func (l *dummyLogger) Error(msg string, fields ...logger.Field) {}

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

func TestDefaultService_ListDrives(t *testing.T) {
	source := new(mockDriveSource)
	mounts := new(mockMountProvider)
	svc := NewDefaultService(source, mounts, &dummyLogger{})

	ctx := context.Background()

	t.Run("list all drives across mounts", func(t *testing.T) {
		mounts.On("ListMounts", ctx).Return([]mount.MountConfig{
			{Path: "/personal", Type: "onedrive", IdentityID: "user1"},
			{Path: "/work", Type: "onedrive", IdentityID: "user2"},
		}, nil).Once()

		source.On("ListDrives", ctx, "/personal").Return([]fs.Drive{
			{ID: "d1", Name: "Personal Drive"},
		}, nil).Once()
		source.On("ListDrives", ctx, "/work").Return([]fs.Drive{
			{ID: "d2", Name: "Work Drive"},
		}, nil).Once()

		drives, err := svc.ListDrives(ctx, "")
		assert.NoError(t, err)
		assert.Len(t, drives, 2)
		assert.Equal(t, "d1", drives[0].ID)
		assert.Equal(t, "d2", drives[1].ID)
	})

	t.Run("filter by identityID", func(t *testing.T) {
		mounts.On("ListMounts", ctx).Return([]mount.MountConfig{
			{Path: "/personal", Type: "onedrive", IdentityID: "user1"},
			{Path: "/work", Type: "onedrive", IdentityID: "user2"},
		}, nil).Once()

		source.On("ListDrives", ctx, "/work").Return([]fs.Drive{
			{ID: "d2", Name: "Work Drive"},
		}, nil).Once()

		drives, err := svc.ListDrives(ctx, "user2")
		assert.NoError(t, err)
		assert.Len(t, drives, 1)
		assert.Equal(t, "d2", drives[0].ID)
	})

	t.Run("skip mounts that don't support drives", func(t *testing.T) {
		mounts.On("ListMounts", ctx).Return([]mount.MountConfig{
			{Path: "/local", Type: "local"},
			{Path: "/personal", Type: "onedrive", IdentityID: "user1"},
		}, nil).Once()

		source.On("ListDrives", ctx, "/local").Return(nil, assert.AnError).Once()
		source.On("ListDrives", ctx, "/personal").Return([]fs.Drive{
			{ID: "d1", Name: "Personal"},
		}, nil).Once()

		drives, err := svc.ListDrives(ctx, "")
		assert.NoError(t, err)
		assert.Len(t, drives, 1)
		assert.Equal(t, "d1", drives[0].ID)
	})
}

func TestDefaultService_ResolveDrive(t *testing.T) {
	source := new(mockDriveSource)
	mounts := new(mockMountProvider)
	svc := NewDefaultService(source, mounts, &dummyLogger{})

	ctx := context.Background()

	mounts.On("ListMounts", ctx).Return([]mount.MountConfig{
		{Path: "/od", Type: "onedrive", IdentityID: "user1"},
	}, nil)

	source.On("ListDrives", ctx, "/od").Return([]fs.Drive{
		{ID: "d1", Name: "Personal Drive", Type: "personal"},
		{ID: "d2", Name: "Work Drive", Type: "business"},
	}, nil)

	t.Run("resolve by ID", func(t *testing.T) {
		d, err := svc.ResolveDrive(ctx, "d1", "")
		assert.NoError(t, err)
		assert.Equal(t, "d1", d.ID)
	})

	t.Run("resolve by Name", func(t *testing.T) {
		d, err := svc.ResolveDrive(ctx, "Work Drive", "")
		assert.NoError(t, err)
		assert.Equal(t, "d2", d.ID)
	})

	t.Run("resolve case insensitive", func(t *testing.T) {
		d, err := svc.ResolveDrive(ctx, "work drive", "")
		assert.NoError(t, err)
		assert.Equal(t, "d2", d.ID)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := svc.ResolveDrive(ctx, "nonexistent", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestDefaultService_ResolvePersonalDrive(t *testing.T) {
	source := new(mockDriveSource)
	mounts := new(mockMountProvider)
	svc := NewDefaultService(source, mounts, &dummyLogger{})

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mounts.On("ListMounts", ctx).Return([]mount.MountConfig{
			{Path: "/od", Type: "onedrive", IdentityID: "user1"},
		}, nil).Once()

		source.On("ListDrives", ctx, "/od").Return([]fs.Drive{
			{ID: "d1", Name: "Personal", Type: "personal"},
		}, nil).Once()

		d, err := svc.ResolvePersonalDrive(ctx, "")
		assert.NoError(t, err)
		assert.Equal(t, "d1", d.ID)
	})

	t.Run("no personal drive", func(t *testing.T) {
		mounts.On("ListMounts", ctx).Return([]mount.MountConfig{
			{Path: "/od", Type: "onedrive", IdentityID: "user1"},
		}, nil).Once()

		source.On("ListDrives", ctx, "/od").Return([]fs.Drive{
			{ID: "d2", Name: "Business", Type: "business"},
		}, nil).Once()

		_, err := svc.ResolvePersonalDrive(ctx, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no personal drive found")
	})
}
