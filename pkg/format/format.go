package format

import "io"

// Format represents a supported output presentation style.
type Format string

const (
	// FormatTable renders data as a human-readable ASCII table.
	FormatTable Format = "table"
	// FormatJSON renders data as a machine-readable JSON object.
	FormatJSON Format = "json"
	// FormatYAML renders data as a machine-readable YAML document.
	FormatYAML Format = "yaml"
	// FormatValue renders a single raw value, omitting any formatting or headers.
	FormatValue Format = "value"
	// FormatShort renders a simplified view of the data, typically just names or IDs.
	FormatShort Format = "short"
)

// Formatter defines the interface for rendering domain data to an [io.Writer].
type Formatter interface {
	// Format writes the provided data to the writer in the specific presentation style.
	Format(w io.Writer, data any) error
}

// Factory coordinates the creation and selection of [Formatter] instances based on user request.
type Factory interface {
	// Get returns the formatter matching the specified [Format].
	// If the format is unknown, a default formatter (typically [FormatShort]) is returned.
	Get(f Format) Formatter
}
