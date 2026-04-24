package mount

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockConfigRepo struct {
	mock.Mock
}

func (m *mockConfigRepo) GetMounts(ctx context.Context) ([]MountConfig, error) {
	args := m.Called(ctx)
	return args.Get(0).([]MountConfig), args.Error(1)
}

func (m *mockConfigRepo) SaveMounts(ctx context.Context, mounts []MountConfig) error {
	args := m.Called(ctx, mounts)
	return args.Error(0)
}

type mockValidator struct {
	mock.Mock
}

func (m *mockValidator) ValidateOptions(opts map[string]string) error {
	args := m.Called(opts)
	return args.Error(0)
}

func (m *mockValidator) ProvideOptions() []MountOption {
	args := m.Called()
	return args.Get(0).([]MountOption)
}

func TestMountService_ListMounts(t *testing.T) {
	repo := new(mockConfigRepo)
	service := NewMountService(repo)

	expectedMounts := []MountConfig{
		{Path: "/onedrive", Type: "onedrive", IdentityID: "user1"},
	}

	repo.On("GetMounts", mock.Anything).Return(expectedMounts, nil)

	mounts, err := service.ListMounts(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedMounts, mounts)
	repo.AssertExpectations(t)
}

func TestMountService_AddMount(t *testing.T) {
	repo := new(mockConfigRepo)
	service := NewMountService(repo)

	newMount := MountConfig{Path: "/new", Type: "local"}
	existingMounts := []MountConfig{
		{Path: "/onedrive", Type: "onedrive"},
	}

	repo.On("GetMounts", mock.Anything).Return(existingMounts, nil)
	repo.On("SaveMounts", mock.Anything, mock.MatchedBy(func(mounts []MountConfig) bool {
		return len(mounts) == 2 && mounts[1].Path == "/new"
	})).Return(nil)

	err := service.AddMount(context.Background(), newMount)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestMountService_AddMount_Update(t *testing.T) {
	repo := new(mockConfigRepo)
	service := NewMountService(repo)

	updatedMount := MountConfig{Path: "/onedrive", Type: "onedrive", IdentityID: "user2"}
	existingMounts := []MountConfig{
		{Path: "/onedrive", Type: "onedrive", IdentityID: "user1"},
	}

	repo.On("GetMounts", mock.Anything).Return(existingMounts, nil)
	repo.On("SaveMounts", mock.Anything, mock.MatchedBy(func(mounts []MountConfig) bool {
		return len(mounts) == 1 && mounts[0].IdentityID == "user2"
	})).Return(nil)

	err := service.AddMount(context.Background(), updatedMount)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestMountService_RemoveMount(t *testing.T) {
	repo := new(mockConfigRepo)
	service := NewMountService(repo)

	existingMounts := []MountConfig{
		{Path: "/onedrive", Type: "onedrive"},
		{Path: "/local", Type: "local"},
	}

	repo.On("GetMounts", mock.Anything).Return(existingMounts, nil)
	repo.On("SaveMounts", mock.Anything, mock.MatchedBy(func(mounts []MountConfig) bool {
		return len(mounts) == 1 && mounts[0].Path == "/local"
	})).Return(nil)

	err := service.RemoveMount(context.Background(), "/onedrive")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestMountService_RemoveMount_NotFound(t *testing.T) {
	repo := new(mockConfigRepo)
	service := NewMountService(repo)

	repo.On("GetMounts", mock.Anything).Return([]MountConfig{}, nil)

	err := service.RemoveMount(context.Background(), "/nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestMountService_ThreadSafety(t *testing.T) {
	repo := new(mockConfigRepo)
	service := NewMountService(repo)

	const iterations = 100
	var wg sync.WaitGroup
	wg.Add(iterations * 3)

	for i := 0; i < iterations; i++ {
		go func(idx int) {
			defer wg.Done()
			mv := new(mockValidator)
			mv.On("ProvideOptions").Return([]MountOption{})
			service.RegisterValidator(fmt.Sprintf("type-%d", idx), mv)
		}(i)

		go func(idx int) {
			defer wg.Done()
			service.GetMountOptions()
		}(i)

		go func(idx int) {
			defer wg.Done()
			service.GetCompletionProvider(fmt.Sprintf("type-%d", idx))
		}(i)
	}

	wg.Wait()
}
