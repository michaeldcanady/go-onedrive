package ls

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/common/sorting"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
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
func CreateLSCmd(c didomain.Container) *cobra.Command {
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
		Long: `You can list the files and folders within a specific OneDrive path. This
command provides various formatting options like long listing, tree view,
and JSON output. You can also sort and filter the results.`,
		Example: `  # List items in the root of your OneDrive
  odc ls /

  # List items in a specific folder using long format
  odc ls /Documents -l

  # Display a recursive tree view of a folder
  odc ls /Projects --tree`,
		Args: cobra.MaximumNArgs(1),

		PreRunE: func(_ *cobra.Command, args []string) error {
			if long {
				opts.Format = "long"
			}
			if tree {
				opts.Format = "tree"
			}
			if len(args) > 0 {
				opts.Path = args[0]
			}
			return opts.Validate()
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.Stdin = cmd.InOrStdin()
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			lsCmd := NewLsCmd(c)
			return lsCmd.Run(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.IncludeAll, allFlagLon, allFlagShort, false, allFlagUsage)
	cmd.Flags().StringVarP(&opts.Format, formatFlagLong, formatFlagShort, formatFlagDefault, formatFlagUsage)
	cmd.Flags().StringVar(&opts.SortProperty, sortFlagLong, sortFlagDefault, sortFlagUsage)
	cmd.Flags().BoolVar(&opts.FoldersOnly, foldersOnlyFlagLong, false, foldersOnlyFlagUsage)
	cmd.Flags().BoolVar(&opts.FilesOnly, filesOnlyFlagLong, false, filesOnlyFlagUsage)
	cmd.Flags().BoolVarP(&opts.Recursive, recursiveFlagLong, recursiveFlagShort, false, recursiveFlagUsage)

	cmd.Flags().BoolVarP(&long, longFlagLong, longFlagShort, false, longFlagUsage)
	cmd.Flags().BoolVar(&tree, treeFlagLong, false, treeFlagUsage)

	return cmd
}
