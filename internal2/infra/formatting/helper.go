package formatting

import domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"

func displayName(it domainfs.Item) string {
	name := it.Path + "/" + it.Name
	if it.Type == domainfs.ItemTypeFolder {
		return name + "/"
	}
	return name
}
