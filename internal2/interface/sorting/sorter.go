package sorting

import (
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

type Sorter interface {
	Sort(items []domainfs.Item) ([]domainfs.Item, error)
}
