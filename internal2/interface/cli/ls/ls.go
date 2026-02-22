// Package ls provides the command-line interface for listing OneDrive items.
package ls

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/sorting"
	"github.com/spf13/cobra"
)

const (
	allFlagLon   = "all"
	allFlagShort = "a"
	allFlagUsage = "show hidden items (names starting with '.')"

	formatFlagLong    = "format"
	formatFlagShort   = "f"
	formatFlagUsage   = "output format (e.g., json, yaml, long, short, tree)"
	formatFlagDefault = "short"

	longFlagLong  = "long"
	longFlagShort = "l"
	longFlagUsage = "use a long listing format"

	treeFlagLong  = "tree"
	treeFlagUsage = "list contents in a tree-like format"

	loggerID    = "cli"
	commandName = "ls"

	filesOnlyFlagLong  = "files-only"
	filesOnlyFlagUsage = "show only files"

	foldersOnlyFlagLong  = "folders-only"
	foldersOnlyFlagUsage = "show only folders"

	sortFlagLong    = "sort"
	sortFlagUsage   = "sorts files by the specified field (e.g., name, size, modified)"
	sortFlagDefault = "name"

	recursiveFlagLong  = "recursive"
	recursiveFlagShort = "R"
	recursiveFlagUsage = "list subdirectories recursively"
)

var (
	supportedFormats    = []string{"json", "yaml", "yml", "long", "short", "tree"}
	supportedProperties = []string{"name", "size", "modified"}
)

// CreateLSCmd constructs and returns the cobra.Command for the ls operation.
// It initializes flags and sets up the execution logic using LsCmd.
func CreateLSCmd(c di.Container) *cobra.Command {
	opts := Options{
		SortOrder: sorting.DirectionAscending,
	}

	var (
		long bool
		tree bool
	)

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [path]", commandName),
		Short: "List items in a OneDrive path",
		Args:  cobra.MaximumNArgs(1),

		PreRunE: func(_ *cobra.Command, _ []string) error {
			if long {
				opts.Format = "long"
			}
			if tree {
				opts.Format = "tree"
			}
			return opts.Validate()
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Path = args[0]
			}
			lsCmd := NewLsCmd(c)
			return lsCmd.Run(cmd.Context(), cmd, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.IncludeAll, allFlagLon, allFlagShort, false, allFlagUsage)
	cmd.Flags().StringVarP(&opts.Format, formatFlagLong, formatFlagShort, formatFlagDefault, formatFlagUsage)
	cmd.Flags().BoolVarP(&long, longFlagLong, longFlagShort, false, longFlagUsage)
	cmd.Flags().BoolVar(&tree, treeFlagLong, false, treeFlagUsage)
	cmd.Flags().StringVar(&opts.SortProperty, sortFlagLong, sortFlagDefault, sortFlagUsage)
	cmd.Flags().BoolVar(&opts.FoldersOnly, foldersOnlyFlagLong, false, foldersOnlyFlagUsage)
	cmd.Flags().BoolVar(&opts.FilesOnly, filesOnlyFlagLong, false, filesOnlyFlagUsage)
	cmd.Flags().BoolVarP(&opts.Recursive, recursiveFlagLong, recursiveFlagShort, false, recursiveFlagUsage)

	return cmd
}
