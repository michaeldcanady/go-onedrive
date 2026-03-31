package formatting

import (
	"encoding/json"
	"io"
)

// JSONFormatter implements OutputFormatter to render items as indented JSON.
type JSONFormatter struct{}

// NewJSONFormatter initializes a new instance of the JSONFormatter.
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// Format encodes the provided items into a JSON array and writes them to the output stream.
func (f *JSONFormatter) Format(w io.Writer, items []any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(items)
}
