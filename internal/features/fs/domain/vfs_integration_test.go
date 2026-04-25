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

	t.Run("ReadFile with identity", func(t *testing.T) {
		mBackend.On("IdentityProvider").Return("onedrive").Once()
		mID.On("Token", mock.Anything, "onedrive", mock.Anything).Return(&proto.GetTokenResponse{
			Token: &proto.AccessToken{Token: "test-token"},
		}, nil).Once()

		// VFS.getToken marshals the entire proto response to string
		mBackend.On("Open", mock.Anything, mock.MatchedBy(func(s string) bool { return s != "" }), "", "/file.txt").
			Return(io.NopCloser(bytes.NewBufferString("content")), nil).Once()

		reader, err := v.ReadFile(ctx, &fs.URI{Provider: "/od", Path: "/file.txt"}, fs.ReadOptions{})
		assert.NoError(t, err)
		content, _ := io.ReadAll(reader)
		assert.Equal(t, "content", string(content))
		mBackend.AssertExpectations(t)
		mID.AssertExpectations(t)
	})
}

func TestVFS_Integration_CrossBackendCopy(t *testing.T) {
	ctx := context.Background()
	v := NewVFS(nil)

	mSrc := new(mockBackend)
	mDst := new(mockBackend)

	v.Mount("/src", mSrc)
	v.Mount("/dst", mDst)

	mSrc.On("IdentityProvider").Return("").Maybe()
	mDst.On("IdentityProvider").Return("").Maybe()

	// Setup: Src has item, Dst will receive it
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
}
