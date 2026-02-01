package filtering

import (
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

type Filterer interface {
	Filter(items []domainfs.Item) ([]domainfs.Item, error)
}
