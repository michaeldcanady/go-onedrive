package formatting

import (
	"io"

	"gopkg.in/yaml.v3"
)

// YAMLFormatter implements OutputFormatter to render items as YAML documents.
type YAMLFormatter struct{}

// NewYAMLFormatter initializes a new instance of the YAMLFormatter.
func NewYAMLFormatter() *YAMLFormatter {
	return &YAMLFormatter{}
}

// Format marshals the provided items into a YAML document and writes it to the output stream.
func (f *YAMLFormatter) Format(w io.Writer, items []any) error {
	out, err := yaml.Marshal(items)
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	return err
}
