package formatting

import domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"

func displayName(it domainfs.Item) string {
	return it.Name
}
