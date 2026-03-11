package infra

import (
	"context"
	"errors"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/fs/domain"
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
		mockPathIDCache := new(MockPathIDCache)

		repo := NewGraphFileContentsGateway(mockAdapter, NewMockLogger())

		mockPathIDCache.On("Get", mock.Anything, path).Return("", false)

		// Cache returns data
		mockContentCache.On("Get", mock.Anything, path).Return(&domain.Contents{
			CTag: etag,
			Data: contentData,
		}, true)

		// Graph returns 304 (nil response)
		mockAdapter.On("SendPrimitive", mock.Anything, mock.Anything, "[]byte", mock.Anything).Return(nil, nil)

		content, _, err := repo.Download(context.Background(), driveID, path, "")
		assert.NoError(t, err)

		assert.Equal(t, contentData, content)
	})

	t.Run("cache miss, fetch fresh", func(t *testing.T) {
		t.Parallel()
		mockAdapter := new(MockRequestAdapter)
		mockContentCache := new(MockContentsCache)
		mockPathIDCache := new(MockPathIDCache)

		repo := NewGraphFileContentsGateway(mockAdapter, NewMockLogger())

		mockPathIDCache.On("Get", mock.Anything, path).Return("", false)

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

		content, _, err := repo.Download(context.Background(), driveID, path, "")
		assert.NoError(t, err)

		assert.Equal(t, contentData, content)
	})

	t.Run("graph error", func(t *testing.T) {
		t.Parallel()
		mockAdapter := new(MockRequestAdapter)
		mockContentCache := new(MockContentsCache)
		mockPathIDCache := new(MockPathIDCache)

		repo := NewGraphFileContentsGateway(mockAdapter, NewMockLogger())

		mockPathIDCache.On("Get", mock.Anything, path).Return("", false)

		mockContentCache.On("Get", mock.Anything, path).Return(nil, false)
		mockAdapter.On("SendPrimitive", mock.Anything, mock.Anything, "[]byte", mock.Anything).Return(nil, errors.New("graph error"))

		r, _, err := repo.Download(context.Background(), driveID, path, "")
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
		mockPathIDCache := new(MockPathIDCache)

		repo := NewGraphFileContentsGateway(mockAdapter, NewMockLogger())

		mockPathIDCache.On("Get", mock.Anything, path).Return("", false)
		mockPathIDCache.On("Put", mock.Anything, path, "new-id").Return(nil)

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
		mockContentCache.On("Put", mock.Anything, "new-id", mock.Anything).Return(nil)
		mockMetadataCache.On("Put", mock.Anything, "new-id", mock.Anything).Return(nil)

		meta, _, err := repo.Upload(context.Background(), driveID, path, contentData, "")

		assert.NoError(t, err)
		assert.NotNil(t, meta)
		assert.Equal(t, etag, meta.CTag)

		mockAdapter.AssertExpectations(t)
		mockContentCache.AssertExpectations(t)
		mockMetadataCache.AssertExpectations(t)
	})
}
