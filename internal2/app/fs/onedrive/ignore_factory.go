package onedrive

import (
	"github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	infrafs "github.com/michaeldcanady/go-onedrive/internal2/infra/fs"
)

func NewIgnoreMatcherFactory() fs.IgnoreMatcherFactory {
	return infrafs.NewIgnoreMatcherFactory()
}
