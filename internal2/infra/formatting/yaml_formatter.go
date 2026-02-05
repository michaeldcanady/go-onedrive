package formatting

import (
	"io"

	"go.yaml.in/yaml/v3"
)

type YAMLFormatter struct{}

func (f *YAMLFormatter) Format(w io.Writer, items any) error {
	out, err := yaml.Marshal(items)
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	return err
}
