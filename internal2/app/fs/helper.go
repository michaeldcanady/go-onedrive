package fs

import (
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	infrafile "github.com/michaeldcanady/go-onedrive/internal2/infra/file"
)

func mapToFSItem(src *infrafile.DriveItem, path string) domainfs.Item {
	t := domainfs.ItemTypeFile
	if src.IsFolder {
		t = domainfs.ItemTypeFolder
	}
	return domainfs.Item{
		ID:       src.ID,
		Path:     path,
		Name:     src.Name,
		Type:     t,
		Size:     src.Size,
		Modified: src.Modified,
	}
}

func mapToFSItems(src []*infrafile.DriveItem, path string) []domainfs.Item {
	result := make([]domainfs.Item, len(src))
	for i, elem := range src {
		result[i] = mapToFSItem(elem, path)
	}
	return result
}
