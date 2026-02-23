// Package ls provides the command-line interface for listing OneDrive items.
package ls

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/filtering"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/formatting"
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
			start := time.Now()

			ctx := cmd.Context()
			if ctx != nil {
				ctx = context.Background()
			}

			logger, err := util.EnsureLogger(c, loggerID)
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("starting ls command",
				infralogging.String("format", format),
				infralogging.Bool("includeAll", includeAll),
				infralogging.Bool("foldersOnly", foldersOnly),
				infralogging.Bool("filesOnly", filesOnly),
				infralogging.String("sortProperty", sortProperty),
				infralogging.String("sortDirection", sortOrder.String()),
			)

			// Filesystem service
			fsSvc := c.FS()
			if fsSvc == nil {
				logger.Error("filesystem service is nil")
				return util.NewCommandErrorWithNameWithMessage(commandName, "filesystem service is nil")
			}

			// Resolve path
			path := ""
			if len(args) > 0 {
				opts.Path = args[0]
			}
			logger.Debug("path resolved", infralogging.String("path", path))

			item, err := fsSvc.Get(ctx, path)
			if err != nil {
				logger.Error("failed to get item", infralogging.String("error", err.Error()))
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			items := []fs.Item{item}
			if item.Type == fs.ItemTypeFolder {
				logger.Debug("listing items from filesystem")
				items, err = fsSvc.List(ctx, path, domainfs.ListOptions{
					Recursive: recursive,
				})
				if err != nil {
					logger.Error("failed to list items", infralogging.String("error", err.Error()))
					return util.NewCommandErrorWithNameWithError(commandName, err)
				}
				logger.Info("items retrieved", infralogging.Int("count", len(items)))
			}

			// Filtering
			logger.Debug("initializing filterer")
			filterer, err := filtering.NewFilterFactory().Create(filterOpts...)
			if err != nil {
				logger.Error("failed to initialize filterer", infralogging.String("error", err.Error()))
				return util.NewCommandError(commandName, "failed to initialize filter", err)
			}

			logger.Debug("applying filters")
			items, err = filterer.Filter(items)
			if err != nil {
				logger.Error("failed to filter items", infralogging.String("error", err.Error()))
				return util.NewCommandError(commandName, "failed to filter items", err)
			}
			logger.Info("items after filtering", infralogging.Int("count", len(items)))

			if format != "tree" {
				// Sorting
				logger.Debug("initializing sorter")
				sorter, err := sorting.NewSorterFactory().Create(sortOpts...)
				if err != nil {
					logger.Error("failed to initialize sorter", infralogging.String("error", err.Error()))
					return util.NewCommandError(commandName, "failed to initialize sorter", err)
				}

				logger.Debug("sorting items")
				items, err = sorter.Sort(items)
				if err != nil {
					logger.Error("failed to sort items", infralogging.String("error", err.Error()))
					return util.NewCommandError(commandName, "failed to sort items", err)
				}
			}

			// Formatting
			logger.Debug("initializing formatter", infralogging.String("format", format))
			formatter, err := formatting.NewFormatterFactory().Create(format)
			if err != nil {
				logger.Error("failed to initialize formatter", infralogging.String("error", err.Error()))
				return util.NewCommandError(commandName, "failed to initialize formatter", err)
			}

			logger.Debug("formatting output")
			if err := formatter.Format(cmd.OutOrStdout(), items); err != nil {
				logger.Error("failed to format items", infralogging.String("error", err.Error()))
				return util.NewCommandError(commandName, "failed to format items", err)
			}

			logger.Info("ls command completed",
				infralogging.Duration("duration", time.Since(start)),
				infralogging.Int("finalItemCount", len(items)),
			)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&includeAll, allFlagLon, allFlagShort, false, allFlagUsage)
	cmd.Flags().StringVarP(&format, formatFlagLong, formatFlagShort, formatFlagDefault, formatFlagUsage)
	cmd.Flags().StringVar(&sortProperty, sortFlagLong, sortFlagDefault, sortFlagUsage)
	cmd.Flags().BoolVar(&foldersOnly, foldersOnlyFlagLong, false, foldersOnlyFlagUsage)
	cmd.Flags().BoolVar(&filesOnly, filesOnlyFlagLong, false, filesOnlyFlagUsage)
	cmd.Flags().BoolVarP(&recursive, recursiveFlagLong, recursiveFlagShort, false, recursiveFlagUsage)

	return cmd
}
