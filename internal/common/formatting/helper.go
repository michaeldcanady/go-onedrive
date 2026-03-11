package formatting

import domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/shared/domain"

func displayName(it domainfs.Item) string {
	return it.Name
}
