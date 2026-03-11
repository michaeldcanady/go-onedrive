package formatting

import (
	"encoding/json"
	"io"

	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

type JSONFormatter struct{}

func (f *JSONFormatter) Format(w io.Writer, items []domainfs.Item) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(items)
}
