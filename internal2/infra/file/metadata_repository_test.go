package file

import (
	"context"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMetadataRepository_GetByPath(t *testing.T) {
	t.Parallel()

	driveID := "drive-1"
	path := "/test.txt"
	etag := "etag-1"

	t.Run("cache hit", func(t *testing.T) {
		t.Parallel()
		mockAdapter := new(MockRequestAdapter)
		mockMetadataCache := new(MockMetadataCache)
		mockListingCache := new(MockListingCache)

		repo := NewMetadataRepository(mockAdapter, mockMetadataCache, mockListingCache)

		// Cache hit
		cached := &file.Metadata{ID: "id-1", Name: "test.txt", ETag: etag}
		mockMetadataCache.On("Get", mock.Anything, path).Return(cached, true)

		// Graph returns 304 (nil response)
		mockAdapter.On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

		got, err := repo.GetByPath(context.Background(), driveID, path, file.MetadataGetOptions{})
		assert.NoError(t, err)
		assert.Equal(t, cached, got)
	})

	t.Run("cache miss, fetch fresh", func(t *testing.T) {
		t.Parallel()
		mockAdapter := new(MockRequestAdapter)
		mockMetadataCache := new(MockMetadataCache)
		mockListingCache := new(MockListingCache)

		repo := NewMetadataRepository(mockAdapter, mockMetadataCache, mockListingCache)

		// Cache miss
		mockMetadataCache.On("Get", mock.Anything, path).Return(nil, false)

		// Graph returns fresh item
		mockItem := models.NewDriveItem()
		id := "id-1"
		name := "test.txt"
		mockItem.SetId(&id)
		mockItem.SetName(&name)
		mockItem.SetETag(&etag)
		mockItem.SetFolder(models.NewFolder())

		mockAdapter.On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockItem, nil)

		// Store in cache
		mockMetadataCache.On("Put", mock.Anything, mock.Anything).Return(nil)

		got, err := repo.GetByPath(context.Background(), driveID, path, file.MetadataGetOptions{})
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, etag, got.ETag)
	})
}

func TestMetadataRepository_ListByPath(t *testing.T) {
	t.Parallel()

	driveID := "drive-1"
	path := "/folder"
	etag := "etag-folder"

	t.Run("listing cache hit", func(t *testing.T) {
		t.Parallel()
		mockAdapter := new(MockRequestAdapter)
		mockMetadataCache := new(MockMetadataCache)
		mockListingCache := new(MockListingCache)

		repo := NewMetadataRepository(mockAdapter, mockMetadataCache, mockListingCache)

		// 1. GetByPath for parent (cache hit, 304)
		parent := &file.Metadata{ID: "id-folder", Name: "folder", ETag: etag}
		mockMetadataCache.On("Get", mock.Anything, path).Return(parent, true)
		mockAdapter.On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

		// 2. Listing cache hit
		listing := &Listing{ETag: etag, ChildIDs: []string{"id-child-1"}}
		mockListingCache.On("Get", mock.Anything, path).Return(listing, true)

		// 3. Resolve child IDs
		child := &file.Metadata{ID: "id-child-1", Name: "child.txt"}
		mockMetadataCache.On("Get", mock.Anything, "id-child-1").Return(child, true)

		got, err := repo.ListByPath(context.Background(), driveID, path, file.MetadataListOptions{})
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, child, got[0])
	})

	t.Run("cache miss, fetch fresh children", func(t *testing.T) {
		t.Parallel()
		mockAdapter := new(MockRequestAdapter)
		mockMetadataCache := new(MockMetadataCache)
		mockListingCache := new(MockListingCache)

		repo := NewMetadataRepository(mockAdapter, mockMetadataCache, mockListingCache)

		// 1. GetByPath for parent (cache miss)
		mockMetadataCache.On("Get", mock.Anything, path).Return(nil, false)
		parentItem := models.NewDriveItem()
		pid := "id-folder"
		pname := "folder"
		parentItem.SetId(&pid)
		parentItem.SetName(&pname)
		parentItem.SetETag(&etag)
		parentItem.SetFolder(models.NewFolder())
		mockAdapter.On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(parentItem, nil).Once()
		mockMetadataCache.On("Put", mock.Anything, mock.Anything).Return(nil)

		// 2. Listing cache miss
		mockListingCache.On("Get", mock.Anything, path).Return(nil, false)

		// 3. Fetch children from Graph
		childItem := models.NewDriveItem()
		cid := "id-child-1"
		cname := "child.txt"
		childItem.SetId(&cid)
		childItem.SetName(&cname)
		childItem.SetFile(models.NewFile())

		coll := models.NewDriveItemCollectionResponse()
		coll.SetValue([]models.DriveItemable{childItem})
		mockAdapter.On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(coll, nil).Once()

		// 4. Store children and listing in cache
		mockMetadataCache.On("Put", mock.Anything, mock.Anything).Return(nil)
		mockListingCache.On("Put", mock.Anything, path, mock.Anything).Return(nil)

		got, err := repo.ListByPath(context.Background(), driveID, path, file.MetadataListOptions{})
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, cid, got[0].ID)
	})
}
