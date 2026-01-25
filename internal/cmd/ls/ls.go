package ls

import (
	"cmp"
	"context"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
)

const (
	longLongFlag  = "long"
	longShortFlag = "l"
	longUsage     = "use long listing format"

	allLongFlag  = "all"
	allShortFlag = "a"
	allUsage     = "show hidden items (names starting with '.')"

	formatLongFlag  = "format"
	formatShortFlag = "f"
	formatUsage     = "output format: json|yaml"
)

func CreateLSCmd(c *di.Container1, logger logging.Logger) *cobra.Command {
	var (
		long   bool
		all    bool
		format string
	)

	cmd := &cobra.Command{
		Use:   "ls [path]",
		Short: "List items in a OneDrive path",
		Long: `List files and folders stored in your OneDrive.

By default, hidden items (names beginning with '.') are not shown.
Use --all to include them, or --long for a detailed listing.`,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			// 🔥 Lazy-load the DriveService here
			drive, err := c.DriveService()
			if err != nil {
				return fmt.Errorf("failed to initialize drive service: %w", err)
			}

			path := ""
			if len(args) == 1 {
				path = args[0]
			}

			logger.Debug("ls command invoked",
				logging.String("path", path),
				logging.Bool("long", long),
				logging.Bool("all", all),
				logging.String("format", format),
			)

			// 🔥 Use the drive service directly
			items, err := collectItems(ctx, drive, path, logger)
			if err != nil {
				return err
			}

			if !all {
				items = filterHidden(items)
			}

			sortItems(items)

			structured := make([]Item, len(items))
			for i, it := range items {
				structured[i] = toItem(it)
			}

			switch format {
			case "":
				if long {
					printLong(structured)
				} else {
					printShort(structured)
				}
			case "json":
				return printJSON(structured)
			case "yaml", "yml":
				return printYAML(structured)
			default:
				return fmt.Errorf("invalid output format: %s (expected json|yaml)", format)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&long, longLongFlag, longShortFlag, false, longUsage)
	cmd.Flags().BoolVarP(&all, allLongFlag, allShortFlag, false, allUsage)
	cmd.Flags().StringVarP(&format, formatLongFlag, formatShortFlag, "", formatUsage)

	return cmd
}

func collectItems(
	ctx context.Context,
	iter driveChildIterator,
	path string,
	logger logging.Logger,
) ([]models.DriveItemable, error) {

	var (
		items   []models.DriveItemable
		iterErr error
	)

	iter.ChildrenIterator(ctx, path)(func(item models.DriveItemable, err error) bool {
		if err != nil {
			logger.Error("iterator returned error", logging.Any("error", err))
			iterErr = err
			return false
		}
		items = append(items, item)
		return true
	})

	if iterErr != nil {
		logger.Error("failed to iterate children", logging.Any("error", iterErr))
		return nil, iterErr
	}

	logger.Debug("items retrieved", logging.Int("count", len(items)))
	return items, nil
}

func filterHidden(items []models.DriveItemable) []models.DriveItemable {
	out := items[:0]
	for _, it := range items {
		if !isHidden(safeName(it)) {
			out = append(out, it)
		}
	}
	return out
}

func sortItems(items []models.DriveItemable) {
	slices.SortFunc(items, func(a, b models.DriveItemable) int {
		return cmp.Compare(safeName(a), safeName(b))
	})
}

func safeName(item models.DriveItemable) string {
	name := ""
	if item.GetName() != nil {
		name = *item.GetName()
	}
	if item.GetFolder() != nil {
		name += "/"
	}
	return name
}

func isHidden(name string) bool {
	return len(name) > 0 && name[0] == '.'
}

func detectTerminalWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80
	}
	return w
}

func getSize(it models.DriveItemable) int64 {
	if it.GetFolder() != nil {
		return 0
	}
	if it.GetSize() != nil {
		return *it.GetSize()
	}
	return 0
}

func getModifiedTime(it models.DriveItemable) time.Time {
	if it.GetLastModifiedDateTime() != nil {
		return *it.GetLastModifiedDateTime()
	}
	return time.Time{}
}
