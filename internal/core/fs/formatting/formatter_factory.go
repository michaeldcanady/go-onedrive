package formatting

import (
	"fmt"
)

// FormatterFactory provides operations for initializing configured formatter instances.
type FormatterFactory struct{}

// NewFormatterFactory initializes a new instance of the FormatterFactory.
func NewFormatterFactory() *FormatterFactory {
	return &FormatterFactory{}
}

// Create returns an OutputFormatter corresponding to the requested format name (e.g., "json", "tree", "long").
// If no format is specified, it defaults to "short".
func (f *FormatterFactory) Create(format string) (OutputFormatter, error) {
	if format == "" {
		format = "short"
	}

	switch format {
	case "short":
		return NewHumanShortFormatter(Terminal{}), nil
	case "long":
		return NewHumanLongFormatter(), nil
	case "json":
		return NewJSONFormatter(), nil
	case "yaml", "yml":
		return NewYAMLFormatter(), nil
	case "tree":
		return NewTreeFormatter(), nil
	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}
}
