package app

import (
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

func convertMetadataToFSType(ype domain.ItemType) domainfs.ItemType {
	switch ype {
	case domain.ItemTypeFile:
		return domainfs.ItemTypeFile
	case domain.ItemTypeFolder:
		return domainfs.ItemTypeFolder
	}
	return domainfs.ItemTypeUnknown
}

func convertMetadataToItem(src *domain.Metadata) domainfs.Item {
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
