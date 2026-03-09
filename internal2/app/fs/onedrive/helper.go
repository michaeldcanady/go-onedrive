package onedrive

import (
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

func convertMetadataToFSType(ype file.ItemType) domainfs.ItemType {
	switch ype {
	case file.ItemTypeFile:
		return domainfs.ItemTypeFile
	case file.ItemTypeFolder:
		return domainfs.ItemTypeFolder
	}
	return domainfs.ItemTypeUnknown
}

func convertMetadataToItem(src *file.Metadata) domainfs.Item {
	if src == nil {
		return domainfs.Item{}
	}

	var modified time.Time
	if src.ModifiedAt != nil {
		modified = *src.ModifiedAt
	}

	return domainfs.Item{
		ID:       src.ID,
		Path:     src.Path,
		Name:     src.Name,
		Size:     src.Size,
		Type:     convertMetadataToFSType(src.Type),
		Modified: modified,
		ETag:     src.ETag,
	}
}
