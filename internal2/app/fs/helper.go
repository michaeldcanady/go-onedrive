package fs

import (
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	infrafile "github.com/michaeldcanady/go-onedrive/internal2/infra/file"
)

func mapToFSItem(src *infrafile.DriveItem) domainfs.Item {
	t := domainfs.ItemTypeFile
	if src.IsFolder {
		t = domainfs.ItemTypeFolder
	}

	return domainfs.Item{
		ID:       src.ID,
		Path:     src.PathWithoutDrive,
		Name:     src.Name,
		Type:     t,
		Size:     src.Size,
		Modified: src.Modified,
	}
}

func mapToFSItems(src []*infrafile.DriveItem) []domainfs.Item {
	result := make([]domainfs.Item, len(src))
	for i, elem := range src {
		result[i] = mapToFSItem(elem)
	}
	return result
}

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
	}
}
