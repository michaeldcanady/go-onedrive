package fs

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVFS_Functional_UserWorkflow(t *testing.T) {
	ctx := context.Background()
	
	// Initialize the functional unit (VFS + URIFactory)
	vfs := NewVFS(nil)
	factory := NewURIFactory(vfs)
	
	// Setup backends
	mRemote := new(mockBackend)
	mLocal := new(mockBackend)
	
	vfs.Mount("/od", mRemote)
	vfs.Mount("/local", mLocal)
	
	t.Run("list and then cat workflow", func(t *testing.T) {
		// 1. User wants to list /od/docs
		mRemote.On("IdentityProvider").Return("").Maybe()
		mRemote.On("List", mock.Anything, "", "", "/docs").Return([]fs.Item{
			{Name: "resume.pdf", Type: fs.TypeFile},
		}, nil).Once()
		
		uri, err := factory.FromString("/od/docs")
		assert.NoError(t, err)
		
		items, err := vfs.List(ctx, uri, fs.ListOptions{})
		assert.NoError(t, err)
		assert.Len(t, items, 1)
		assert.Equal(t, "resume.pdf", items[0].Name)
		
		// 2. User wants to read /od/docs/resume.pdf
		mRemote.On("Open", mock.Anything, "", "", "/docs/resume.pdf").
			Return(io.NopCloser(bytes.NewBufferString("pdf data")), nil).Once()
			
		fileUri, err := factory.FromString("/od/docs/resume.pdf")
		assert.NoError(t, err)
		
		reader, err := vfs.ReadFile(ctx, fileUri, fs.ReadOptions{})
		assert.NoError(t, err)
		content, _ := io.ReadAll(reader)
		assert.Equal(t, "pdf data", string(content))
		
		mRemote.AssertExpectations(t)
	})

	t.Run("cross-mount copy workflow", func(t *testing.T) {
		// User wants to copy /od/file.txt to /local/backup.txt
		mRemote.On("Open", mock.Anything, "", "", "/file.txt").
			Return(io.NopCloser(bytes.NewBufferString("important data")), nil).Once()
		
		mLocal.On("Create", mock.Anything, "", "", "/backup.txt", mock.Anything).
			Return(fs.Item{Name: "backup.txt"}, nil).Once()
			
		srcUri, _ := factory.FromString("/od/file.txt")
		dstUri, _ := factory.FromString("/local/backup.txt")
		
		err := vfs.Copy(ctx, srcUri, dstUri, fs.CopyOptions{})
		assert.NoError(t, err)
		
		mRemote.AssertExpectations(t)
		mLocal.AssertExpectations(t)
	})
}
