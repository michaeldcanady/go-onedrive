package infra

import (
	"testing"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stretchr/testify/assert"
)

func TestMapItemToMetadata(t *testing.T) {
	t.Parallel()

	now := time.Now()
	id := "item-id"
	name := "test-item"
	etag := "etag-123"
	ctag := "ctag-456"
	size := int64(1024)
	mimeType := "text/plain"

	tests := []struct {
		name     string
		input    models.DriveItemable
		expected *domain.Metadata
	}{
		{
			name: "folder item",
			input: func() models.DriveItemable {
				item := models.NewDriveItem()
				item.SetId(&id)
				item.SetName(&name)
				item.SetSize(&size)
				item.SetETag(&etag)
				item.SetCTag(&ctag)
				item.SetCreatedDateTime(&now)
				item.SetLastModifiedDateTime(&now)
				item.SetFolder(models.NewFolder())

				parent := models.NewItemReference()
				pid := "parent-id"
				ppath := "drive-id:/path/to/parent"
				parent.SetId(&pid)
				parent.SetPath(&ppath)
				item.SetParentReference(parent)

				return item
			}(),
			expected: &domain.Metadata{
				ID:         id,
				Name:       name,
				FullPath:   "drive-id:/path/to/parent",
				Path:       "/path/to/parent",
				Size:       size,
				ETag:       etag,
				CTag:       ctag,
				ParentID:   "parent-id",
				CreatedAt:  &now,
				ModifiedAt: &now,
				Type:       domain.ItemTypeFolder,
			},
		},
		{
			name: "file item",
			input: func() models.DriveItemable {
				item := models.NewDriveItem()
				item.SetId(&id)
				item.SetName(&name)
				item.SetSize(&size)
				item.SetETag(&etag)
				item.SetCTag(&ctag)
				item.SetCreatedDateTime(&now)
				item.SetLastModifiedDateTime(&now)

				f := models.NewFile()
				f.SetMimeType(&mimeType)
				item.SetFile(f)

				parent := models.NewItemReference()
				pid := "parent-id"
				ppath := "drive-id:/path/to/parent"
				parent.SetId(&pid)
				parent.SetPath(&ppath)
				item.SetParentReference(parent)

				return item
			}(),
			expected: &domain.Metadata{
				ID:         id,
				Name:       name,
				FullPath:   "drive-id:/path/to/parent",
				Path:       "/path/to/parent",
				Size:       size,
				MimeType:   mimeType,
				ETag:       etag,
				CTag:       ctag,
				ParentID:   "parent-id",
				CreatedAt:  &now,
				ModifiedAt: &now,
				Type:       domain.ItemTypeFile,
			},
		},
		{
			name: "root item",
			input: func() models.DriveItemable {
				item := models.NewDriveItem()
				item.SetId(&id)
				item.SetName(&name)
				item.SetFolder(models.NewFolder())
				// No parent reference for root usually, or it's minimal
				return item
			}(),
			expected: &domain.Metadata{
				ID:       id,
				Name:     name,
				FullPath: name,
				Path:     name,
				Type:     domain.ItemTypeFolder,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := mapItemToMetadata(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}
