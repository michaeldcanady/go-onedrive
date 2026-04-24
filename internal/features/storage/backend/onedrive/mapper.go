package onedrive

import (
	"time"

	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func mapItemToSharedItem(it models.DriveItemable, itemPath string) fs.Item {
	if it == nil {
		return fs.Item{}
	}

	id := ""
	if it.GetId() != nil {
		id = *it.GetId()
	}

	name := ""
	if it.GetName() != nil {
		name = *it.GetName()
	}

	size := int64(0)
	if it.GetSize() != nil {
		size = *it.GetSize()
	}

	itemType := fs.TypeFolder
	if it.GetFile() != nil {
		itemType = fs.TypeFile
	}

	modifiedAt := it.GetLastModifiedDateTime()
	etag := ""
	if it.GetETag() != nil {
		etag = *it.GetETag()
	}

	var mTime time.Time
	if modifiedAt != nil {
		mTime = *modifiedAt
	}

	return fs.Item{
		ID:         id,
		Name:       name,
		Path:       itemPath,
		Type:       itemType,
		Size:       size,
		ModifiedAt: mTime,
		ETag:       etag,
	}
}
