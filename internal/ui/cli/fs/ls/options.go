package ls

import (
	"fmt"
	"io"
	"slices"
	"strings"

	fs "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	formatting "github.com/michaeldcanady/go-onedrive/pkg/format"
)

// Options provides the user-facing settings for the drive ls command.
type Options struct {
	// Path is the filesystem path to list.
	Path string
	// URI is the parsed and resolved filesystem location.
	URI *fs.URI
	// Recursive determines whether to list subdirectories.
	Recursive bool

	// Format is the output format (e.g., "short", "long", "json", "tree").
	Format formatting.Format
	// SortFields is the list of fields to sort by (e.g., "name", "size", "modified").
	SortFields []string
	// SortDescending determines whether to sort in descending order.
	SortDescending bool
	// All determines whether to include hidden files.
	All bool
	// Stdout is the destination for standard output messages.
	Stdout io.Writer
}

// Validate ensures that the provided options are consistent and valid.
func (o *Options) Validate() error {
	// Validate SortFields
	validSortFields := []string{"name", "size", "modified"}
	for _, field := range o.SortFields {
		if !slices.Contains(validSortFields, strings.ToLower(field)) {
			return fmt.Errorf("invalid sorting field '%s'; please use one of the following valid fields: %s",
				field, strings.Join(validSortFields, ", "))
		}
	}

	// Validate Format
	// Note: FormatUnknown is set when the provided format string doesn't match known types.
	if o.Format == formatting.FormatUnknown {
		validFormats := []string{"short", "long", "json", "yaml", "tree", "table"}
		return fmt.Errorf("unknown output format specified; please provide a valid format such as: %s",
			strings.Join(validFormats, ", "))
	}

	// Cross-field validation: Recursive mode restrictions
	if o.Recursive {
		allowedRecursiveFormats := []formatting.Format{
			formatting.FormatTree,
			formatting.FormatLong,
			formatting.FormatJSON,
			formatting.FormatYAML,
		}

		if !slices.Contains(allowedRecursiveFormats, o.Format) {
			var allowedStrings []string
			for _, f := range allowedRecursiveFormats {
				allowedStrings = append(allowedStrings, f.String())
			}

			return fmt.Errorf("recursive mode (-r/--recursive) is not supported with the '%s' format; "+
				"please use a compatible format like: %s",
				o.Format.String(), strings.Join(allowedStrings, ", "))
		}
	}

	return nil
}
