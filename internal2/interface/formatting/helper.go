package formatting

import domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"

func displayName(it domainfs.Item) string {
	return it.Name
}
