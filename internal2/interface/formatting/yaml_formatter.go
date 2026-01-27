package formatting

import (
	"io"

	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	"go.yaml.in/yaml/v3"
)

type YAMLFormatter struct{}

func (f *YAMLFormatter) Format(w io.Writer, items []domainfs.Item) error {
	out, err := yaml.Marshal(items)
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	return err
}
