package formatting

// FormatterFactory provides operations for initializing configured formatter instances.
type FormatterFactory struct{}

// NewFormatterFactory initializes a new instance of the FormatterFactory.
func NewFormatterFactory() *FormatterFactory {
	return &FormatterFactory{}
}

// Create returns an OutputFormatter corresponding to the requested format name (e.g., "json", "tree", "long").
// If no format is specified, it defaults to "short".
func (f *FormatterFactory) Create(format Format) (OutputFormatter, error) {
	if format == FormatUnknown {
		format = FormatShort
	}

	switch format {
	case FormatShort:
		return NewHumanShortFormatter(nil), nil
	case FormatLong:
		return NewHumanLongFormatter(), nil
	case FormatJSON:
		return NewJSONFormatter(), nil
	case FormatYAML:
		return NewYAMLFormatter(), nil
	case FormatTree:
		return NewTreeFormatter(), nil
	default:
		return nil, NewUnsupportedFormatError(format, nil)
	}
}
