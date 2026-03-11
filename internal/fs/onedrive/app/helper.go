package app

import (
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/fs/shared/domain"
)

func convertMetadataToFSType(ype domain.ItemType) domain.ItemType {
	switch ype {
	case domain.ItemTypeFile:
		return domain.ItemTypeFile
	case domain.ItemTypeFolder:
		return domain.ItemTypeFolder
	}
	return domain.ItemTypeUnknown
}

func convertMetadataToItem(src *domain.Metadata) domain.Item {
	if src == nil {
		return domain.Item{}
	}

	var modified time.Time
	if src.ModifiedAt != nil {
		modified = *src.ModifiedAt
	}

	return domain.Item{
		ID:       src.ID,
		Path:     src.Path,
		Name:     src.Name,
		Size:     src.Size,
		Type:     convertMetadataToFSType(src.Type),
		Modified: modified,
		ETag:     src.ETag,
	}
}
