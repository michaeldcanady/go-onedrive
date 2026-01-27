package ls

import (
	"context"
	"fmt"
	"time"

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
		format string
		all    bool
	)

	cmd := &cobra.Command{
		Use:   "ls [path]",
		Short: "List items in a OneDrive path",
		Args:  cobra.MaximumNArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			start := time.Now()

			// Add base context fields for this command
			ctx = logging.WithFields(ctx,
				logging.String("command", "ls"),
				logging.String("profile", c.Options.ProfileName),
			)

			// Get base logger and enrich it with context
			baseLogger, _ := c.LoggerService.GetLogger("cli")
			logger := baseLogger.WithContext(ctx)

			logger.Info("ls command started")

			// Resolve path argument
			path := ""
			if len(args) == 1 {
				path = args[0]
			}

			// Add path to context
			ctx = logging.WithFields(ctx, logging.String("path", path))
			logger = logger.WithContext(ctx)

			logger.Debug("initializing services")

			// Load services
			fs, err := c.FileSystemService(ctx)
			if err != nil {
				logger.Error("failed to initialize filesystem service", logging.Any("error", err))
				return fmt.Errorf("failed to initialize filesystem service: %w", err)
			}
			logger.Info("initialized filesystem service")

			ds, err := c.DriveService(ctx)
			if err != nil {
				logger.Error("failed to initialize drive service", logging.Any("error", err))
				return fmt.Errorf("failed to initialize drive service: %w", err)
			}
			logger.Info("initialized drive service")

			// Resolve drive
			drive, err := ds.ResolvePersonalDrive(ctx)
			if err != nil {
				logger.Error("failed to resolve personal drive", logging.Any("error", err))
				return fmt.Errorf("failed to resolve personal drive: %w", err)
			}
			if drive == nil {
				logger.Error("no working drive selected")
				return fmt.Errorf("no working drive selected")
			}
			logger.Info("drive selected", logging.String("driveID", drive.ID))

			// Add drive ID to context
			ctx = logging.WithFields(ctx, logging.String("driveID", drive.ID))
			logger = logger.WithContext(ctx)

			logger.Debug("resolving item")

			// Resolve item
			item, err := fs.ResolveItem(ctx, drive.ID, path)
			if err != nil {
				logger.Warn("failed to resolve item", logging.Any("error", err))
				return handleDomainError("ls", path, err)
			}
			logger.Info("resolved item",
				logging.String("itemID", item.ID),
				logging.Bool("isFolder", item.IsFolder),
			)

			var items []*driveservice.DriveItem

			if !item.IsFolder {
				items = []*driveservice.DriveItem{item}
			} else {
				logger.Debug("listing children")
				items, err = fs.ListChildren(ctx, drive.ID, path)
				if err != nil {
					logger.Warn("failed to list children", logging.Any("error", err))
					return handleDomainError("ls", path, err)
				}
				logger.Info("found children", logging.Int("count", len(items)))
			}

			if !all {
				before := len(items)
				logger.Debug("filtering hidden items")

				items = filterHiddenDomain(items)

				logger.Info("filtered hidden items",
					logging.Int("before", before),
					logging.Int("after", len(items)),
					logging.Int("removed", before-len(items)),
				)
			}

			logger.Debug("sorting items")
			sortDomainItems(items)

			// Formatter selection
			logger.Info("selecting formatter", logging.String("format", format))
			formatter, err := NewFormatterFactory().Create(format)
			if err != nil {
				logger.Error("invalid format", logging.String("format", format))
				return err
			}

			logger.Info("rendering output", logging.String("format", format))
			if err := formatter.Format(items); err != nil {
				logger.Error("failed to render output", logging.Any("error", err))
				return err
			}

			logger.Info("ls command completed",
				logging.Duration("elapsed", time.Since(start)),
			)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&all, allFlagLong, allFlagShort, false, allFlagUsage)
	cmd.Flags().StringVarP(&format, formatLongFlag, formatShortFlag, "", formatUsage)

	return cmd
}
