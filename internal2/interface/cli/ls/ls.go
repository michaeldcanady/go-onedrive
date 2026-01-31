package ls

import (
	"fmt"
	"os"

	applogging "github.com/michaeldcanady/go-onedrive/internal2/app/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/filtering"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/formatting"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/sorting"
	"github.com/spf13/cobra"
)

const (
	allFlagLong  = "all"
	allFlagShort = "a"
	allFlagUsage = "show hidden items (names starting with '.')"

	formatLongFlag  = "format"
	formatShortFlag = "f"
	formatUsage     = "output format: json|yaml|long|short"

	loggerID    = "cli"
	commandName = "ls"
)

func CreateLSCmd(c di.Container) *cobra.Command {
	var (
		format       string
		includeAll   bool
		foldersOnly  bool
		filesOnly    bool
		sortProperty = "Name"
		sortOrder    = sorting.DirectionAscending
		sortOpts     = []sorting.SortingOption{sorting.WithDirection(sortOrder)}
		filterOpts   = []filtering.FilterOption{}
	)

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [path]", commandName),
		Short: "List items in a OneDrive path",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if foldersOnly && filesOnly {
				return NewCommandErrorWithNameWithMessage(commandName, "can't use --folders-only and --files-only together")
			}

			if includeAll {
				filterOpts = append(filterOpts, filtering.IncludeAll())
			} else {
				filterOpts = append(filterOpts, filtering.ExcludeAll())
			}

			if filesOnly {
				filterOpts = append(filterOpts, filtering.WithItemType(domainfs.ItemTypeFile))
			} else if foldersOnly {
				filterOpts = append(filterOpts, filtering.WithItemType(domainfs.ItemTypeFolder))
			}

			if sortProperty != "" {
				sortOpts = append(sortOpts, sorting.WithField(sortProperty))
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, err := ensureLogger(c)
			if err != nil {
				return NewCommandErrorWithNameWithError(commandName, err)
			}

			path := ""
			if len(args) > 0 {
				path = args[0]
			}
			logger.Debug("path resolved", infralogging.String("path", path))

			fsSvc := c.FS()
			if fsSvc == nil {
				return NewCommandErrorWithNameWithMessage(commandName, "filesystem service is nil")
			}

			items, err := fsSvc.List(cmd.Context(), path, domainfs.ListOptions{})
			if err != nil {
				return NewCommandErrorWithNameWithError(commandName, err)
			}

			// Filtering

			filterer, err := filtering.NewFilterFactory().Create(filterOpts...)
			if err != nil {
				return NewCommandError(commandName, "failed to initialize filter", err)
			}

			items, err = filterer.Filter(items)
			if err != nil {
				return NewCommandError(commandName, "failed to filter items", err)
			}

			// Sorting
			sorter, err := sorting.NewSorterFactory().Create()
			if err != nil {
				return NewCommandError(commandName, "failed to initialize sorter", err)
			}

			items, err = sorter.Sort(items)
			if err != nil {
				return NewCommandError(commandName, "failed to sort items", err)
			}

			// Formatting
			formatter, err := formatting.NewFormatterFactory().Create(format)
			if err != nil {
				return NewCommandError(commandName, "failed to initialize formatter", err)
			}

			if err := formatter.Format(os.Stdout, items); err != nil {
				return NewCommandError(commandName, "failed to format items", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&includeAll, allFlagLong, allFlagShort, false, allFlagUsage)
	cmd.Flags().StringVarP(&format, formatLongFlag, formatShortFlag, "", formatUsage)
	cmd.Flags().BoolVar(&foldersOnly, "folders-only", false, "show only folders")
	cmd.Flags().BoolVar(&filesOnly, "files-only", false, "show only files")

	return cmd
}

// ensureLogger retrieves or creates the CLI logger.
func ensureLogger(c di.Container) (infralogging.Logger, error) {
	logger, err := c.Logger().GetLogger(loggerID)
	if err == applogging.ErrUnknownLogger {
		return c.Logger().CreateLogger(loggerID)
	}
	return logger, err
}
