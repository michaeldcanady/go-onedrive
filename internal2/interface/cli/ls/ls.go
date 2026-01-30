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
		sortType     = "property"
		sortProperty = "Name"
		sortOrder    = sorting.DirectionAscending
	)

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [path]", commandName),
		Short: "List items in a OneDrive path",
		Args:  cobra.MaximumNArgs(1),

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
			filterType := "hidden"
			if includeAll {
				filterType = "none"
			}

			filterer, err := filtering.NewFilterFactory().Create(filterType)
			if err != nil {
				return NewCommandError(commandName, "failed to initialize filter", err)
			}

			items, err = filterer.Filter(items)
			if err != nil {
				return NewCommandError(commandName, "failed to filter items", err)
			}

			// Sorting
			sorter, err := sorting.NewSorterFactory().Create(
				sortType,
				sorting.WithField(sortProperty),
				sorting.WithDirection(sortOrder),
			)
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
