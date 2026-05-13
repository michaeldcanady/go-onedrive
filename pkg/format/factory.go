package format

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

// Tabular is an optional interface that data structures can implement to provide
// structured data for the [FormatTable] formatter.
type Tabular interface {
	// TableHeaders returns the list of column titles.
	TableHeaders() []string

	// TableRows returns the 2D slice of string data for each row.
	TableRows() [][]string
}

type factory struct {
	formatters map[Format]Formatter
}

// NewFactory returns a new [Factory] initialized with standard formatters.
func NewFactory() Factory {
	return &factory{
		formatters: map[Format]Formatter{
			FormatJSON:  &jsonFormatter{},
			FormatYAML:  &yamlFormatter{},
			FormatTable: &tableFormatter{},
			FormatValue: &valueFormatter{},
			FormatShort: &shortFormatter{},
		},
	}
}

func (f *factory) Get(format Format) Formatter {
	if formatter, ok := f.formatters[format]; ok {
		return formatter
	}
	return f.formatters[FormatShort]
}

type jsonFormatter struct{}

func (f *jsonFormatter) Format(w io.Writer, data any) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

type yamlFormatter struct{}

func (f *yamlFormatter) Format(w io.Writer, data any) error {
	return yaml.NewEncoder(w).Encode(data)
}

type tableFormatter struct{}

func (f *tableFormatter) Format(w io.Writer, data any) error {
	tabular, ok := data.(Tabular)
	if !ok {
		// Fallback to short formatter if not tabular
		return (&shortFormatter{}).Format(w, data)
	}

	rows := tabular.TableRows()
	if len(rows) == 0 {
		fmt.Fprintln(w, "No items found.")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	headers := tabular.TableHeaders()
	for i, h := range headers {
		fmt.Fprint(tw, h)
		if i < len(headers)-1 {
			fmt.Fprint(tw, "\t")
		}
	}
	fmt.Fprintln(tw)

	for _, row := range rows {
		for i, cell := range row {
			fmt.Fprint(tw, cell)
			if i < len(row)-1 {
				fmt.Fprint(tw, "\t")
			}
		}
		fmt.Fprintln(tw)
	}

	return tw.Flush()
}

type shortFormatter struct{}

func (f *shortFormatter) Format(w io.Writer, data any) error {
	switch v := data.(type) {
	case []string:
		for _, item := range v {
			fmt.Fprintln(w, item)
		}
		return nil
	case string:
		fmt.Fprintln(w, v)
		return nil
	default:
		// If it's a slice of something else, try to print its string representation
		// Or if it has a String() method.
		if s, ok := v.(fmt.Stringer); ok {
			fmt.Fprintln(w, s.String())
			return nil
		}
		fmt.Fprintln(w, v)
		return nil
	}
}

type valueFormatter struct{}

func (f *valueFormatter) Format(w io.Writer, data any) error {
	if s, ok := data.(string); ok {
		_, err := fmt.Fprintln(w, s)
		return err
	}
	_, err := fmt.Fprintln(w, data)
	return err
}
