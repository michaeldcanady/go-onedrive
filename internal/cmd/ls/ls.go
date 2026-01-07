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

	"github.com/michaeldcanady/go-onedrive/internal/logging"
)

const (
	longLongArg  = "long"
	longShortArg = "l"
	longArgUsage = "use long listing format"

	allLongArg  = "all"
	allShortArg = "a"
	allArgUsage = "show hidden items (names starting with '.')"
)

func CreateLSCmd(iter driveChildIterator, logger logging.Logger) *cobra.Command {
	var long bool
	var all bool
	var format string

	lsCmd := &cobra.Command{
		Use:          "ls [path]",
		Short:        "List items in a OneDrive path",
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			path := ""
			if len(args) == 1 {
				path = args[0]
			}

			logger.Debug("ls command invoked",
				logging.String("path", path),
				logging.Any("long", long),
				logging.Any("all", all),
			)

			var items []models.DriveItemable
			var cmdErr error

			iter.ChildrenIterator(ctx, path)(func(item models.DriveItemable, err error) bool {
				if err != nil {
					logger.Error("iterator returned error", logging.Any("error", err))
					cmdErr = err
					return false
				}
				items = append(items, item)
				return true
			})
			if cmdErr != nil {
				logger.Error("failed to iterate children", logging.Any("error", cmdErr))
				return cmdErr
			}

			logger.Debug("items retrieved", logging.Int("count", len(items)))

			// Filter hidden items unless --all is set
			if !all {
				before := len(items)
				filtered := items[:0]
				for _, it := range items {
					if !isHidden(safeName(it)) {
						filtered = append(filtered, it)
					}
				}
				items = filtered
				logger.Debug("filtered hidden items",
					logging.Int("before", before),
					logging.Int("after", len(items)),
				)
			}

			// Sort by name
			slices.SortFunc(items, func(a, b models.DriveItemable) int {
				return cmp.Compare(safeName(a), safeName(b))
			})

			// Convert to LSItem
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

	lsCmd.Flags().BoolVarP(&long, longLongArg, longShortArg, false, longArgUsage)
	lsCmd.Flags().BoolVarP(&all, allLongArg, allShortArg, false, allArgUsage)
	lsCmd.Flags().StringVarP(&format, "format", "f", "", "output format: json|yaml")

	return lsCmd
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
