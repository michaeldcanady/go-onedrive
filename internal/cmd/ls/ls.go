package ls

import (
	"fmt"

	"github.com/spf13/cobra"

	driveservice "github.com/michaeldcanady/go-onedrive/internal/app/file_service"
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
)

const (
	allFlagLong  = "all"
	allFlagShort = "a"
	allFlagUsage = "show hidden items (names starting with '.')"

	formatLongFlag  = "format"
	formatShortFlag = "f"
	formatUsage     = "output format: json|yaml|long|short"
)

func CreateLSCmd(c *di.Container) *cobra.Command {
	var (
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

			// Load services
			fs, err := c.FileSystemService(ctx)
			if err != nil {
				return fmt.Errorf("failed to initialize filesystem service: %w", err)
			}

			ds, err := c.DriveService(ctx)
			if err != nil {
				return fmt.Errorf("failed to initialize drive service: %w", err)
			}

			// Resolve working drive (for now: personal drive)
			drive, err := ds.ResolvePersonalDrive(ctx)
			if err != nil {
				return fmt.Errorf("failed to resolve personal drive: %w", err)
			}
			if drive == nil {
				return fmt.Errorf("no working drive selected")
			}

			// Determine path
			path := ""
			if len(args) == 1 {
				path = args[0]
			}

			logger.Debug("ls invoked",
				logging.String("drive", drive.ID),
				logging.String("path", path),
				logging.String("format", format),
				logging.Bool("all", all),
			)

			// Resolve item
			item, err := fs.ResolveItem(ctx, drive.ID, path)
			if err != nil {
				return handleDomainError("ls", path, err)
			}

			var items []*driveservice.DriveItem

			if !item.IsFolder {
				items = []*driveservice.DriveItem{item}
			} else {
				items, err = fs.ListChildren(ctx, drive.ID, path)
				if err != nil {
					return handleDomainError("ls", path, err)
				}
			}

			// Filter hidden
			if !all {
				items = filterHiddenDomain(items)
			}

			// Sort
			sortDomainItems(items)

			formatter, err := NewFormatterFactory().Create(format)
			if err != nil {
				return err
			}

			return formatter.Format(items)
		},
	}

	cmd.Flags().BoolVarP(&all, allFlagLong, allFlagShort, false, allFlagUsage)
	cmd.Flags().StringVarP(&format, formatLongFlag, formatShortFlag, "", formatUsage)

	return cmd
}
