package file

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContentsRepository_Download(t *testing.T) {
	t.Parallel()

	driveID := "drive-1"
	path := "/test.txt"
	contentData := []byte("hello world")
	etag := "etag-1"

	t.Run("cache hit", func(t *testing.T) {
		t.Parallel()
		mockAdapter := new(MockRequestAdapter)
		mockContentCache := new(MockContentsCache)
		mockMetadataCache := new(MockMetadataCache)

		repo := NewContentsRepository(mockAdapter, mockContentCache, mockMetadataCache)

		// Cache returns data
		mockContentCache.On("Get", mock.Anything, path).Return(&file.Contents{
			CTag: etag,
			Data: contentData,
		}, true)

		// Graph returns 304 (nil response)
		mockAdapter.On("SendPrimitive", mock.Anything, mock.Anything, "[]byte", mock.Anything).Return(nil, nil)

		r, err := repo.Download(context.Background(), driveID, path, file.DownloadOptions{})
		assert.NoError(t, err)
		defer r.Close()

		gotData, _ := io.ReadAll(r)
		assert.Equal(t, contentData, gotData)
	})

	t.Run("cache miss, fetch fresh", func(t *testing.T) {
		t.Parallel()
		mockAdapter := new(MockRequestAdapter)
		mockContentCache := new(MockContentsCache)
		mockMetadataCache := new(MockMetadataCache)

		repo := NewContentsRepository(mockAdapter, mockContentCache, mockMetadataCache)

		// Cache miss
		mockContentCache.On("Get", mock.Anything, path).Return(nil, false)

		// Graph returns data
		mockAdapter.On("SendPrimitive", mock.Anything, mock.Anything, "[]byte", mock.Anything).Return(contentData, nil)

		// Should cache the result if headers are present (mocking headers is hard here without response info)
		// The current implementation relies on headerOpt to extract CTag.
		// Since we can't easily inject the response headers via the mock adapter return values (it returns `any`),
		// we might not trigger the `Put` logic unless we mock the middleware or response handler.
		// However, looking at the code: `headerOpt.GetResponseHeaders()` is used.
		// `SendPrimitive` returns `(any, error)`. It doesn't modify the `headerOpt` directly in a way we can control easily here
		// unless `SendPrimitive` logic in the real adapter populates it.
		// For now, we'll assume no cache Put if we can't simulate headers easily, or just verify the download.

		r, err := repo.Download(context.Background(), driveID, path, file.DownloadOptions{})
		assert.NoError(t, err)
		defer r.Close()

		gotData, _ := io.ReadAll(r)
		assert.Equal(t, contentData, gotData)
	})

	t.Run("graph error", func(t *testing.T) {
		t.Parallel()
		mockAdapter := new(MockRequestAdapter)
		mockContentCache := new(MockContentsCache)
		mockMetadataCache := new(MockMetadataCache)

		repo := NewContentsRepository(mockAdapter, mockContentCache, mockMetadataCache)

		mockContentCache.On("Get", mock.Anything, path).Return(nil, false)
		mockAdapter.On("SendPrimitive", mock.Anything, mock.Anything, "[]byte", mock.Anything).Return(nil, errors.New("graph error"))

		r, err := repo.Download(context.Background(), driveID, path, file.DownloadOptions{})
		assert.Error(t, err)
		assert.Nil(t, r)
	})
}

func TestContentsRepository_Upload(t *testing.T) {
	t.Parallel()

	driveID := "drive-1"
	path := "/test.txt"
	contentData := []byte("new content")
	etag := "etag-new"

	t.Run("successful upload", func(t *testing.T) {
		t.Parallel()
		mockAdapter := new(MockRequestAdapter)
		mockContentCache := new(MockContentsCache)
		mockMetadataCache := new(MockMetadataCache)

		repo := NewContentsRepository(mockAdapter, mockContentCache, mockMetadataCache)

		// 1. Force is false, so check cache for If-Match
		mockContentCache.On("Get", mock.Anything, path).Return(nil, false)

		// 2. Mock Graph Put
		mockItem := models.NewDriveItem()
		id := "new-id"
		name := "test.txt"
		mockItem.SetId(&id)
		mockItem.SetName(&name)
		mockItem.SetCTag(&etag)
		f := models.NewFile()
		mt := "text/plain"
		f.SetMimeType(&mt)
		mockItem.SetFile(f)

		mockAdapter.On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockItem, nil)

		// 3. Mock Cache Puts
		mockContentCache.On("Put", mock.Anything, path, mock.Anything).Return(nil)
		mockMetadataCache.On("Put", mock.Anything, mock.Anything).Return(nil)

		body := bytes.NewReader(contentData)
		meta, err := repo.Upload(context.Background(), driveID, path, body, file.UploadOptions{})

		assert.NoError(t, err)
		assert.NotNil(t, meta)
		assert.Equal(t, etag, meta.CTag)

		mockAdapter.AssertExpectations(t)
		mockContentCache.AssertExpectations(t)
		mockMetadataCache.AssertExpectations(t)
	})
}

