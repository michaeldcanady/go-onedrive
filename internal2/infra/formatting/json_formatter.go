package formatting

import (
	"encoding/json"
	"io"
)

type JSONFormatter struct{}

func (f *JSONFormatter) Format(w io.Writer, items any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(items)
}
