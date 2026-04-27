package fs

import (
	"context"
	"io"

	proto "github.com/michaeldcanady/go-onedrive/internal/features/identity/proto"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/stretchr/testify/mock"
)

type mockBackend struct {
	mock.Mock
}

func (m *mockBackend) Name() string             { return m.Called().String(0) }
func (m *mockBackend) IdentityProvider() string { return m.Called().String(0) }
func (m *mockBackend) Stat(ctx context.Context, token, driveID, path string) (fs.Item, error) {
	args := m.Called(ctx, token, driveID, path)
	return args.Get(0).(fs.Item), args.Error(1)
}
func (m *mockBackend) List(ctx context.Context, token, driveID, path string) ([]fs.Item, error) {
	args := m.Called(ctx, token, driveID, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]fs.Item), args.Error(1)
}
func (m *mockBackend) Open(ctx context.Context, token, driveID, path string) (io.ReadCloser, error) {
	args := m.Called(ctx, token, driveID, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}
func (m *mockBackend) Create(ctx context.Context, token, driveID, path string, r io.Reader) (fs.Item, error) {
	args := m.Called(ctx, token, driveID, path, r)
	return args.Get(0).(fs.Item), args.Error(1)
}
func (m *mockBackend) Mkdir(ctx context.Context, token, driveID, path string) error {
	return m.Called(ctx, token, driveID, path).Error(0)
}
func (m *mockBackend) Remove(ctx context.Context, token, driveID, path string) error {
	return m.Called(ctx, token, driveID, path).Error(0)
}
func (m *mockBackend) Capabilities() fs.Capabilities {
	return m.Called().Get(0).(fs.Capabilities)
}
func (m *mockBackend) Move(ctx context.Context, token, driveID, src, dst string) error {
	return m.Called(ctx, token, driveID, src, dst).Error(0)
}
func (m *mockBackend) Copy(ctx context.Context, token, driveID, src, dst string) error {
	return m.Called(ctx, token, driveID, src, dst).Error(0)
}

type mockTokenProvider struct {
	mock.Mock
}

func (m *mockTokenProvider) Token(ctx context.Context, provider string, req *proto.GetTokenRequest) (*proto.GetTokenResponse, error) {
	args := m.Called(ctx, provider, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*proto.GetTokenResponse), args.Error(1)
}

type mockVFS struct {
	mock.Mock
}

func (m *mockVFS) Resolve(absPath string) (string, string, error) {
	args := m.Called(absPath)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *mockBackend) ListDrives(ctx context.Context, token string) ([]fs.Drive, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]fs.Drive), args.Error(1)
}

func (m *mockBackend) GetPersonalDrive(ctx context.Context, token string) (fs.Drive, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(fs.Drive), args.Error(1)
}
