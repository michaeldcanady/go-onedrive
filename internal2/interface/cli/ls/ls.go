package ls

import (
	"os"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
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
)

func CreateLSCmd(c di.Container) *cobra.Command {
	var (
		format string
		all    bool

		sortType     string            = "property"
		sortProperty string            = "Name"
		sortOrder    sorting.Direction = sorting.DirectionAscending
	)

	cmd := &cobra.Command{
		Use:   "ls [path]",
		Short: "List items in a OneDrive path",
		Args:  cobra.MaximumNArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			path := ""
			if len(args) > 0 {
				path = args[0]
			}

			filesystemService := c.FS()

			fsItems, err := filesystemService.List(cmd.Context(), path, domainfs.ListOptions{})
			if err != nil {
				return err
			}

			sortingOpts := []sorting.SortingOption{sorting.WithDirection(sortOrder), sorting.WithField(sortProperty)}

			sorter, err := sorting.NewSorterFactory().Create(sortType, sortingOpts...)
			if err != nil {
				return err
			}
			fsItems, err = sorter.Sort(fsItems)
			if err != nil {
				return err
			}

			formatter, err := formatting.NewFormatterFactory().Create(format)
			if err != nil {
				return err
			}

			return formatter.Format(os.Stdout, fsItems)
		},
	}

	cmd.Flags().BoolVarP(&all, allFlagLong, allFlagShort, false, allFlagUsage)
	cmd.Flags().StringVarP(&format, formatLongFlag, formatShortFlag, "", formatUsage)

	return cmd
}
