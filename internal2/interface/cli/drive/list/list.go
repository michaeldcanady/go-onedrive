package list

import (
	"context"
	"reflect"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	domaindrive "github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/formatting"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/spf13/cobra"
)

const (
	commandName = "list"
	loggerID    = "cli"
)

func CreateListCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists available drives",
		Args:  cobra.ExactArgs(0),

		PreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
				cmd.SetContext(ctx)
			}

			logger, err := util.EnsureLogger(ctx, container, loggerID)
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("command pre-run started",
				logging.String("command", commandName),
			)

			logger.Debug("registering table formatter",
				logging.String("event", "register_formatter"),
			)

			tableFormatter := formatting.NewTableFormatter(
				formatting.DriveIDColumn,
				formatting.DriveNameColumn,
				formatting.DriveOwnerColumn,
				formatting.DriveReadOnlyColumn,
				formatting.DriveTypeColumn,
			)

			if err := container.Format().RegisterWithType(
				"table",
				reflect.TypeFor[[]*domaindrive.Drive](),
				tableFormatter,
			); err != nil {

				logger.Warn("failed to register table formatter",
					logging.String("event", "register_formatter"),
					logging.Error(err),
				)

				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("command pre-run completed",
				logging.String("command", commandName),
			)

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
				cmd.SetContext(ctx)
			}

			logger, err := util.EnsureLogger(ctx, container, loggerID)
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("command started",
				logging.String("command", commandName),
			)

			logger.Debug("retrieving drives",
				logging.String("event", "list_drives"),
			)

			drives, err := container.Drive().ListDrives(ctx)
			if err != nil {
				logger.Warn("failed to retrieve drives",
					logging.String("event", "list_drives"),
					logging.Error(err),
				)
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("drives retrieved",
				logging.String("event", "list_drives"),
				logging.Int("count", len(drives)),
			)

			logger.Debug("formatting output",
				logging.String("event", "format_output"),
				logging.String("format", "table"),
			)

			if err := container.Format().Format(cmd.OutOrStdout(), "table", drives); err != nil {
				logger.Warn("failed to format output",
					logging.String("event", "format_output"),
					logging.Error(err),
				)
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("command completed",
				logging.String("command", commandName),
				logging.Int("drive_count", len(drives)),
			)

			return nil
		},
	}
}
