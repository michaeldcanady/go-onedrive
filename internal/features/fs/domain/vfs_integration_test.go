package fs

import (
	"bytes"
	"context"
	"io"
	"testing"

	proto "github.com/michaeldcanady/go-onedrive/internal/features/identity/proto"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVFS_Integration_Operations(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ctx context.Context, v *VFS, mID *mockTokenProvider, mB *mockBackend)
		action  func(ctx context.Context, v *VFS) error
		verify  func(t *testing.T, err error)
	}{
		{
			name: "Stat delegation",
			setup: func(ctx context.Context, v *VFS, mID *mockTokenProvider, mB *mockBackend) {
				mB.On("IdentityProvider").Return("").Once()
				mB.On("Stat", mock.Anything, "", "", "/file.txt").Return(fs.Item{Name: "file.txt"}, nil).Once()
			},
			action: func(ctx context.Context, v *VFS) error {
				item, err := v.Stat(ctx, &fs.URI{Provider: "/od", Path: "/file.txt"})
				assert.Equal(t, "file.txt", item.Name)
				return err
			},
			verify: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "List delegation",
			setup: func(ctx context.Context, v *VFS, mID *mockTokenProvider, mB *mockBackend) {
				mB.On("IdentityProvider").Return("").Once()
				mB.On("List", mock.Anything, "", "", "/dir").Return([]fs.Item{{Name: "f1"}}, nil).Once()
			},
			action: func(ctx context.Context, v *VFS) error {
				items, err := v.List(ctx, &fs.URI{Provider: "/od", Path: "/dir"}, fs.ListOptions{})
				assert.Len(t, items, 1)
				return err
			},
			verify: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "ReadFile with identity",
			setup: func(ctx context.Context, v *VFS, mID *mockTokenProvider, mB *mockBackend) {
				mB.On("IdentityProvider").Return("onedrive").Once()
				mID.On("Token", mock.Anything, "onedrive", mock.Anything).Return(&proto.GetTokenResponse{
					Token: &proto.AccessToken{Token: "test-token"},
				}, nil).Once()
				mB.On("Open", mock.Anything, mock.MatchedBy(func(s string) bool { return s != "" }), "", "/file.txt").
					Return(io.NopCloser(bytes.NewBufferString("content")), nil).Once()
			},
			action: func(ctx context.Context, v *VFS) error {
				reader, err := v.ReadFile(ctx, &fs.URI{Provider: "/od", Path: "/file.txt"}, fs.ReadOptions{})
				if err == nil {
					content, _ := io.ReadAll(reader)
					assert.Equal(t, "content", string(content))
				}
				return err
			},
			verify: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "WriteFile delegation",
			setup: func(ctx context.Context, v *VFS, mID *mockTokenProvider, mB *mockBackend) {
				mB.On("IdentityProvider").Return("").Once()
				mB.On("Create", mock.Anything, "", "", "/new.txt", mock.Anything).Return(fs.Item{Name: "new.txt"}, nil).Once()
			},
			action: func(ctx context.Context, v *VFS) error {
				item, err := v.WriteFile(ctx, &fs.URI{Provider: "/od", Path: "/new.txt"}, bytes.NewBufferString("data"), fs.WriteOptions{})
				assert.Equal(t, "new.txt", item.Name)
				return err
			},
			verify: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "Mkdir delegation",
			setup: func(ctx context.Context, v *VFS, mID *mockTokenProvider, mB *mockBackend) {
				mB.On("IdentityProvider").Return("").Once()
				mB.On("Mkdir", mock.Anything, "", "", "/newdir").Return(nil).Once()
			},
			action: func(ctx context.Context, v *VFS) error {
				return v.Mkdir(ctx, &fs.URI{Provider: "/od", Path: "/newdir"})
			},
			verify: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "Remove delegation",
			setup: func(ctx context.Context, v *VFS, mID *mockTokenProvider, mB *mockBackend) {
				mB.On("IdentityProvider").Return("").Once()
				mB.On("Remove", mock.Anything, "", "", "/file.txt").Return(nil).Once()
			},
			action: func(ctx context.Context, v *VFS) error {
				return v.Remove(ctx, &fs.URI{Provider: "/od", Path: "/file.txt"})
			},
			verify: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mID := new(mockTokenProvider)
			v := NewVFS(mID)
			mB := new(mockBackend)
			v.Mount("/od", mB)

			tt.setup(ctx, v, mID, mB)
			err := tt.action(ctx, v)
			tt.verify(t, err)

			mB.AssertExpectations(t)
			mID.AssertExpectations(t)
		})
	}
}

func TestVFS_Integration_CopyMove(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ctx context.Context, v *VFS, mSrc, mDst *mockBackend)
		action  func(ctx context.Context, v *VFS) error
		wantErr bool
	}{
		{
			name: "Cross-backend Copy (fallback)",
			setup: func(ctx context.Context, v *VFS, mSrc, mDst *mockBackend) {
				mSrc.On("IdentityProvider").Return("").Maybe()
				mDst.On("IdentityProvider").Return("").Maybe()
				mSrc.On("Open", mock.Anything, "", "", "/file.txt").
					Return(io.NopCloser(bytes.NewBufferString("data")), nil).Once()
				mDst.On("Create", mock.Anything, "", "", "/file.txt", mock.Anything).
					Return(fs.Item{}, nil).Once()
			},
			action: func(ctx context.Context, v *VFS) error {
				return v.Copy(ctx, 
					&fs.URI{Provider: "/src", Path: "/file.txt"}, 
					&fs.URI{Provider: "/dst", Path: "/file.txt"}, 
					fs.CopyOptions{})
			},
			wantErr: false,
		},
		{
			name: "Same-backend Native Copy",
			setup: func(ctx context.Context, v *VFS, mSrc, mDst *mockBackend) {
				mSrc.On("IdentityProvider").Return("").Maybe()
				mSrc.On("Capabilities").Return(fs.Capabilities{CanCopy: true}).Once()
				mSrc.On("Copy", mock.Anything, "", "", "/file1.txt", "/file2.txt").Return(nil).Once()
			},
			action: func(ctx context.Context, v *VFS) error {
				return v.Copy(ctx, 
					&fs.URI{Provider: "/src", Path: "/file1.txt"}, 
					&fs.URI{Provider: "/src", Path: "/file2.txt"}, 
					fs.CopyOptions{})
			},
			wantErr: false,
		},
		{
			name: "Same-backend Native Move",
			setup: func(ctx context.Context, v *VFS, mSrc, mDst *mockBackend) {
				mSrc.On("IdentityProvider").Return("").Maybe()
				mSrc.On("Capabilities").Return(fs.Capabilities{CanMove: true}).Once()
				mSrc.On("Move", mock.Anything, "", "", "/old.txt", "/new.txt").Return(nil).Once()
			},
			action: func(ctx context.Context, v *VFS) error {
				return v.Move(ctx, 
					&fs.URI{Provider: "/src", Path: "/old.txt"}, 
					&fs.URI{Provider: "/src", Path: "/new.txt"})
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			v := NewVFS(nil)
			mSrc := new(mockBackend)
			mDst := new(mockBackend)
			v.Mount("/src", mSrc)
			v.Mount("/dst", mDst)

			tt.setup(ctx, v, mSrc, mDst)
			err := tt.action(ctx, v)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mSrc.AssertExpectations(t)
			mDst.AssertExpectations(t)
		})
	}
}

func TestVFS_Integration_DriveDiscovery(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(mB *mockBackend)
		action  func(ctx context.Context, v *VFS) error
		wantErr bool
	}{
		{
			name: "ListDrives success",
			setup: func(mB *mockBackend) {
				mB.On("IdentityProvider").Return("").Once()
				mB.On("ListDrives", mock.Anything, "").Return([]fs.Drive{{ID: "d1"}}, nil).Once()
			},
			action: func(ctx context.Context, v *VFS) error {
				drives, err := v.ListDrives(ctx, "/od")
				assert.Len(t, drives, 1)
				return err
			},
		},
		{
			name: "GetPersonalDrive success",
			setup: func(mB *mockBackend) {
				mB.On("IdentityProvider").Return("").Once()
				mB.On("GetPersonalDrive", mock.Anything, "").Return(fs.Drive{ID: "pd"}, nil).Once()
			},
			action: func(ctx context.Context, v *VFS) error {
				drive, err := v.GetPersonalDrive(ctx, "/od")
				assert.Equal(t, "pd", drive.ID)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			v := NewVFS(nil)
			mB := new(mockBackend)
			v.Mount("/od", mB)

			tt.setup(mB)
			err := tt.action(ctx, v)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mB.AssertExpectations(t)
		})
	}
}
