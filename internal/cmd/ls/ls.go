package ls

import (
	"github.com/michaeldcanady/go-onedrive/internal/di2"
	"github.com/michaeldcanady/go-onedrive/internal/fs"
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

func CreateLSCmd(c *di2.Container) *cobra.Command {
	var (
		format string
		all    bool
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

			fsItems, err := filesystemService.List(cmd.Context(), path, fs.ListOptions{})
			if err != nil {
				return err
			}

			sortDomainItems(fsItems)

			formatter, err := NewFormatterFactory().Create(format)
			if err != nil {
				return err
			}

			return formatter.Format(fsItems)
		},
	}

	cmd.Flags().BoolVarP(&all, allFlagLong, allFlagShort, false, allFlagUsage)
	cmd.Flags().StringVarP(&format, formatLongFlag, formatShortFlag, "", formatUsage)

	return cmd
}
