package formatting

import domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"

func displayName(it domainfs.Item) string {
	if it.Type == domainfs.ItemTypeFolder {
		return it.Name + "/"
	}
	return it.Name
}
