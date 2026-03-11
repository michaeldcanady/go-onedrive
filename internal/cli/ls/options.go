// Package ls provides the command-line interface for listing OneDrive items.
package ls

import (
	"fmt"
	"io"
	"slices"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	"github.com/michaeldcanady/go-onedrive/internal/common/sorting"
)

// Options defines the configuration for the ls command.
// It encapsulates all flags and arguments that control the listing behavior,
// including path resolution, output formatting, filtering, and sorting.
type Options struct {
	// Path is the filesystem path to list. If empty, it defaults to the root of the OneDrive.
	Path string

	// Format specifies the output format.
	// Supported values include "short", "long", "json", "yaml", and "tree".
	Format string

	// IncludeAll determines whether to show hidden files (those starting with a dot).
	IncludeAll bool

	// FoldersOnly restricts the output to include only folder items.
	// Cannot be used simultaneously with FilesOnly.
	FoldersOnly bool

	// FilesOnly restricts the output to include only file items.
	// Cannot be used simultaneously with FoldersOnly.
	FilesOnly bool

	// Recursive determines whether to list subdirectories and their contents recursively.
	Recursive bool

	// IgnoreFile is the path to a .gitignore-style file for filtering.
	IgnoreFile string

	// SortProperty specifies the property used for sorting the results.
	// Supported values include "name", "size", and "modified".
	SortProperty string

	// SortOrder specifies the direction of sorting (ascending or descending).
	SortOrder sorting.Direction

	// Stdin is the input stream for the command.
	Stdin io.Reader
	// Stdout is the output stream for the command.
	Stdout io.Writer
	// Stderr is the error stream for the command.
	Stderr io.Writer
}

// Validate checks if the options are valid and consistent according to the business rules.
// It returns a [util.CommandError] if validation fails, such as when conflicting flags
// are used or unsupported formats/properties are specified.
func (o *Options) Validate() error {
	if o.FoldersOnly && o.FilesOnly {
		return util.NewCommandErrorWithNameWithMessage(commandName, "can't use --folders-only and --files-only together")
	}

	if !slices.Contains(supportedFormats, o.Format) {
		return util.NewCommandErrorWithNameWithMessage(commandName, fmt.Sprintf("unsupported format: %s; only supports: %v", o.Format, supportedFormats))
	}

	if !slices.Contains(supportedProperties, o.SortProperty) {
		return util.NewCommandErrorWithNameWithMessage(commandName, fmt.Sprintf("unsupported property: %s; only supports: %v", o.SortProperty, supportedProperties))
	}

	return nil
}
