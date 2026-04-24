package formatting

import "strings"

type Format int32

const (
	FormatUnknown Format = iota
	FormatJSON
	FormatYAML
	FormatLong
	FormatShort
	FormatTable
	FormatTree
)

// String returns the string representation of the format.
func (dt Format) String() string {
	switch dt {
	case FormatJSON:
		return "json"
	case FormatYAML:
		return "yaml"
	case FormatLong:
		return "long"
	case FormatShort:
		return "short"
	case FormatTable:
		return "table"
	case FormatTree:
		return "tree"
	default:
		return "unknown"
	}
}

// NewFormat converts a string to its corresponding Format.
func NewFormat(s string) Format {
	switch strings.ToLower(s) {
	case "json":
		return FormatJSON
	case "yaml":
		return FormatYAML
	case "long":
		return FormatLong
	case "short":
		return FormatShort
	case "table":
		return FormatTable
	case "tree":
		return FormatTree
	default:
		return FormatUnknown
	}
}
