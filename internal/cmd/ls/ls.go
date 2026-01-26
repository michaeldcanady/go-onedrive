package ls

import (
	"fmt"

	"github.com/spf13/cobra"

	driveservice "github.com/michaeldcanady/go-onedrive/internal/app/drive_service"
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
)

const (
	longLongFlag  = "long"
	longShortFlag = "l"
	longUsage     = "use long listing format"

	allFlagLong  = "all"
	allFlagShort = "a"
	allFlagUsage = "show hidden items (names starting with '.')"

	formatLongFlag  = "format"
	formatShortFlag = "f"
	formatUsage     = "output format: json|yaml"
)

func CreateLSCmd(c *di.Container1) *cobra.Command {
	var (
		long   bool
		all    bool
		format string
	)

	cmd := &cobra.Command{
		Use:   "ls [path]",
		Short: "List items in a OneDrive path",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, _ := c.LoggerService.GetLogger("cli")
			ctx := cmd.Context()

			fs, err := c.FileService(ctx)
			if err != nil {
				return fmt.Errorf("failed to initialize filesystem service: %w", err)
			}

			ds, err := c.DriveService2(ctx)
			if err != nil {
				return fmt.Errorf("failed to initialize driveservice service: %w", err)
			}

			// for now only supports user's personal drive
			drive, err := ds.ResolvePersonalDrive(ctx)
			if err != nil {
				return fmt.Errorf("failed to resolve drive %s: %w", "OneDrive", err)
			}
			if drive == nil {
				return fmt.Errorf("no working drive selected")
			}

			path := ""
			if len(args) == 1 {
				path = args[0]
			}

			logger.Debug("ls invoked",
				logging.String("drive", drive.ID),
				logging.String("path", path),
				logging.Bool("long", long),
				logging.Bool("all", all),
			)

			var items []*driveservice.DriveItem

			item, err := fs.ResolveItem(ctx, drive.ID, path)
			if err != nil {
				return handleDomainError("ls", path, err)
			}

			if !item.IsFolder {
				items = []*driveservice.DriveItem{item}
			} else {
				items, err = fs.ListChildren(ctx, drive.ID, path)
				if err != nil {
					return handleDomainError("ls", path, err)
				}
			}

			if !all {
				logger.Debug("filtering items")
				items = filterHiddenDomain(items)
			}

			sortDomainItems(items)

			switch format {
			case "":
				if long {
					printLongDomain(items)
				} else {
					printShortDomain(items)
				}
			case "json":
				return printJSON(items)
			case "yaml", "yml":
				return printYAML(items)
			default:
				return fmt.Errorf("invalid format: %s", format)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&long, longLongFlag, longShortFlag, false, longUsage)
	cmd.Flags().BoolVarP(&all, allFlagLong, allFlagShort, false, allFlagUsage)
	cmd.Flags().StringVarP(&format, formatLongFlag, formatShortFlag, "", formatUsage)

	return cmd
}
