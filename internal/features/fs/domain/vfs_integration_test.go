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
	ctx := context.Background()
	mID := new(mockTokenProvider)
	v := NewVFS(mID)

	mBackend := new(mockBackend)
	v.Mount("/od", mBackend)

	t.Run("Stat delegation", func(t *testing.T) {
		mBackend.On("IdentityProvider").Return("").Once()
		mBackend.On("Stat", mock.Anything, "", "", "/file.txt").Return(fs.Item{Name: "file.txt"}, nil).Once()

		item, err := v.Stat(ctx, &fs.URI{Provider: "/od", Path: "/file.txt"})
		assert.NoError(t, err)
		assert.Equal(t, "file.txt", item.Name)
		mBackend.AssertExpectations(t)
	})

	t.Run("List delegation", func(t *testing.T) {
		mBackend.On("IdentityProvider").Return("").Once()
		mBackend.On("List", mock.Anything, "", "", "/dir").Return([]fs.Item{{Name: "f1"}}, nil).Once()

		items, err := v.List(ctx, &fs.URI{Provider: "/od", Path: "/dir"}, fs.ListOptions{})
		assert.NoError(t, err)
		assert.Len(t, items, 1)
		mBackend.AssertExpectations(t)
	})

	t.Run("ReadFile with identity", func(t *testing.T) {
		mBackend.On("IdentityProvider").Return("onedrive").Once()
		mID.On("Token", mock.Anything, "onedrive", mock.Anything).Return(&proto.GetTokenResponse{
			Token: &proto.AccessToken{Token: "test-token"},
		}, nil).Once()

		mBackend.On("Open", mock.Anything, mock.MatchedBy(func(s string) bool { return s != "" }), "", "/file.txt").
			Return(io.NopCloser(bytes.NewBufferString("content")), nil).Once()

		reader, err := v.ReadFile(ctx, &fs.URI{Provider: "/od", Path: "/file.txt"}, fs.ReadOptions{})
		assert.NoError(t, err)
		content, _ := io.ReadAll(reader)
		assert.Equal(t, "content", string(content))
		mBackend.AssertExpectations(t)
		mID.AssertExpectations(t)
	})

	t.Run("WriteFile delegation", func(t *testing.T) {
		mBackend.On("IdentityProvider").Return("").Once()
		mBackend.On("Create", mock.Anything, "", "", "/new.txt", mock.Anything).Return(fs.Item{Name: "new.txt"}, nil).Once()

		item, err := v.WriteFile(ctx, &fs.URI{Provider: "/od", Path: "/new.txt"}, bytes.NewBufferString("data"), fs.WriteOptions{})
		assert.NoError(t, err)
		assert.Equal(t, "new.txt", item.Name)
		mBackend.AssertExpectations(t)
	})

	t.Run("Mkdir delegation", func(t *testing.T) {
		mBackend.On("IdentityProvider").Return("").Once()
		mBackend.On("Mkdir", mock.Anything, "", "", "/newdir").Return(nil).Once()

		err := v.Mkdir(ctx, &fs.URI{Provider: "/od", Path: "/newdir"})
		assert.NoError(t, err)
		mBackend.AssertExpectations(t)
	})

	t.Run("Remove delegation", func(t *testing.T) {
		mBackend.On("IdentityProvider").Return("").Once()
		mBackend.On("Remove", mock.Anything, "", "", "/file.txt").Return(nil).Once()

		err := v.Remove(ctx, &fs.URI{Provider: "/od", Path: "/file.txt"})
		assert.NoError(t, err)
		mBackend.AssertExpectations(t)
	})
}

func TestVFS_Integration_CopyMove(t *testing.T) {
	ctx := context.Background()
	v := NewVFS(nil)

	mSrc := new(mockBackend)
	mDst := new(mockBackend)

	v.Mount("/src", mSrc)
	v.Mount("/dst", mDst)

	t.Run("Cross-backend Copy (fallback)", func(t *testing.T) {
		mSrc.On("IdentityProvider").Return("").Maybe()
		mDst.On("IdentityProvider").Return("").Maybe()

		mSrc.On("Open", mock.Anything, "", "", "/file.txt").
			Return(io.NopCloser(bytes.NewBufferString("data")), nil).Once()
		mDst.On("Create", mock.Anything, "", "", "/file.txt", mock.Anything).
			Return(fs.Item{}, nil).Once()

		err := v.Copy(ctx, 
			&fs.URI{Provider: "/src", Path: "/file.txt"}, 
			&fs.URI{Provider: "/dst", Path: "/file.txt"}, 
			fs.CopyOptions{})

		assert.NoError(t, err)
		mSrc.AssertExpectations(t)
		mDst.AssertExpectations(t)
	})

	t.Run("Same-backend Native Copy", func(t *testing.T) {
		mSrc.On("IdentityProvider").Return("").Maybe()
		mSrc.On("Capabilities").Return(fs.Capabilities{CanCopy: true}).Once()
		mSrc.On("Copy", mock.Anything, "", "", "/file1.txt", "/file2.txt").Return(nil).Once()

		err := v.Copy(ctx, 
			&fs.URI{Provider: "/src", Path: "/file1.txt"}, 
			&fs.URI{Provider: "/src", Path: "/file2.txt"}, 
			fs.CopyOptions{})

		assert.NoError(t, err)
		mSrc.AssertExpectations(t)
	})

	t.Run("Same-backend Native Move", func(t *testing.T) {
		mSrc.On("IdentityProvider").Return("").Maybe()
		mSrc.On("Capabilities").Return(fs.Capabilities{CanMove: true}).Once()
		mSrc.On("Move", mock.Anything, "", "", "/old.txt", "/new.txt").Return(nil).Once()

		err := v.Move(ctx, 
			&fs.URI{Provider: "/src", Path: "/old.txt"}, 
			&fs.URI{Provider: "/src", Path: "/new.txt"})

		assert.NoError(t, err)
		mSrc.AssertExpectations(t)
	})
}

func TestVFS_Integration_DriveDiscovery(t *testing.T) {
	ctx := context.Background()
	v := NewVFS(nil)

	mBackend := new(mockBackend)
	v.Mount("/od", mBackend)

	t.Run("ListDrives success", func(t *testing.T) {
		mBackend.On("IdentityProvider").Return("").Once()
		mBackend.On("ListDrives", mock.Anything, "").Return([]fs.Drive{{ID: "d1"}}, nil).Once()

		drives, err := v.ListDrives(ctx, "/od")
		assert.NoError(t, err)
		assert.Len(t, drives, 1)
		mBackend.AssertExpectations(t)
	})

	t.Run("GetPersonalDrive success", func(t *testing.T) {
		mBackend.On("IdentityProvider").Return("").Once()
		mBackend.On("GetPersonalDrive", mock.Anything, "").Return(fs.Drive{ID: "pd"}, nil).Once()

		drive, err := v.GetPersonalDrive(ctx, "/od")
		assert.NoError(t, err)
		assert.Equal(t, "pd", drive.ID)
		mBackend.AssertExpectations(t)
	})
}
