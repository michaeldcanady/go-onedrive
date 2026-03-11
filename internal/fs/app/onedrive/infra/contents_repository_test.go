package infra

import (
	"context"
	"errors"
	"testing"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContentsRepository_Download(t *testing.T) {
	t.Parallel()

	driveID := "drive-1"
	path := "/test.txt"
	contentData := []byte("hello world")

	t.Run("fresh fetch", func(t *testing.T) {
		t.Parallel()
		mockAdapter := new(MockRequestAdapter)

		repo := NewGraphFileContentsGateway(mockAdapter, NewMockLogger())

		// Graph returns data
		mockAdapter.On("SendPrimitive", mock.Anything, mock.Anything, "[]byte", mock.Anything).Return(contentData, nil)

		content, _, err := repo.Download(context.Background(), driveID, path, "")
		assert.NoError(t, err)
		assert.Equal(t, contentData, content)
	})

	t.Run("graph error", func(t *testing.T) {
		t.Parallel()
		mockAdapter := new(MockRequestAdapter)

		repo := NewGraphFileContentsGateway(mockAdapter, NewMockLogger())

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

		repo := NewGraphFileContentsGateway(mockAdapter, NewMockLogger())

		// Mock Graph Put
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

		meta, _, err := repo.Upload(context.Background(), driveID, path, contentData, "")

		assert.NoError(t, err)
		assert.NotNil(t, meta)
		assert.Equal(t, etag, meta.CTag)

		mockAdapter.AssertExpectations(t)
	})
}
