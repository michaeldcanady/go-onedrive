package infra

import (
	"context"
	"testing"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMetadataRepository_GetByPath(t *testing.T) {
	t.Parallel()

	driveID := "drive-1"
	path := "/test.txt"
	etag := "etag-1"

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockAdapter := new(MockRequestAdapter)

		repo := NewGraphMetadataGateway(mockAdapter, NewMockLogger())

		// Graph returns fresh item
		mockItem := models.NewDriveItem()
		id := "id-1"
		name := "test.txt"
		mockItem.SetId(&id)
		mockItem.SetName(&name)
		mockItem.SetETag(&etag)
		mockItem.SetFolder(models.NewFolder())

		mockAdapter.On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockItem, nil)

		got, err := repo.GetByPath(context.Background(), driveID, path, "")
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, etag, got.ETag)
	})

	t.Run("304 not modified", func(t *testing.T) {
		t.Parallel()
		mockAdapter := new(MockRequestAdapter)

		repo := NewGraphMetadataGateway(mockAdapter, NewMockLogger())

		// Graph returns 304 (nil response)
		mockAdapter.On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

		got, err := repo.GetByPath(context.Background(), driveID, path, "etag-1")
		assert.NoError(t, err)
		assert.Nil(t, got)
	})
}

func TestMetadataRepository_ListByPath(t *testing.T) {
	t.Parallel()

	driveID := "drive-1"
	path := "/folder"

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockAdapter := new(MockRequestAdapter)

		repo := NewGraphMetadataGateway(mockAdapter, NewMockLogger())

		// Fetch children from Graph
		childItem := models.NewDriveItem()
		cid := "id-child-1"
		cname := "child.txt"
		childItem.SetId(&cid)
		childItem.SetName(&cname)
		childItem.SetFile(models.NewFile())

		coll := models.NewDriveItemCollectionResponse()
		coll.SetValue([]models.DriveItemable{childItem})
		mockAdapter.On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(coll, nil).Once()

		got, err := repo.ListByPath(context.Background(), driveID, path, "")
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, cid, got[0].ID)
	})

	t.Run("304 not modified", func(t *testing.T) {
		t.Parallel()
		mockAdapter := new(MockRequestAdapter)

		repo := NewGraphMetadataGateway(mockAdapter, NewMockLogger())

		// Graph returns 304 (nil response)
		mockAdapter.On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()

		got, err := repo.ListByPath(context.Background(), driveID, path, "etag-folder")
		assert.NoError(t, err)
		assert.Nil(t, got)
	})
}
