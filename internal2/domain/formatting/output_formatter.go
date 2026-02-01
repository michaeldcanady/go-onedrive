package formatting

import (
	"io"

	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

type OutputFormatter interface {
	Format(w io.Writer, items []domainfs.Item) error
}
