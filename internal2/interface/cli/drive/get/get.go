package get

import (
	"context"
	"reflect"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	domaindrive "github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/formatting"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/spf13/cobra"
)

const (
	commandName = "get"
	loggerID    = "cli"
)

func CreateGetCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get information for a drive",
		Args:  cobra.ExactArgs(1),

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
				reflect.TypeOf([]*domaindrive.Drive{}),
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
				logging.Strings("args", args),
			)

			id := strings.ToLower(strings.TrimSpace(args[0]))
			if id == "" {
				logger.Warn("empty drive id provided",
					logging.String("event", "validate_input"),
				)
				return util.NewCommandErrorWithNameWithMessage(commandName, "id is empty")
			}

			logger.Debug("resolving drive",
				logging.String("event", "resolve_drive"),
				logging.String("drive_id", id),
			)

			drive, err := container.Drive().ResolveDrive(ctx, id)
			if err != nil {
				logger.Warn("failed to resolve drive",
					logging.String("event", "resolve_drive"),
					logging.Error(err),
					logging.String("drive_id", id),
				)
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Debug("formatting drive output",
				logging.String("event", "format_output"),
				logging.String("drive_id", drive.ID),
			)

			if err := container.Format().Format(
				cmd.OutOrStdout(),
				"table",
				[]*domaindrive.Drive{drive},
			); err != nil {

				logger.Warn("failed to format output",
					logging.String("event", "format_output"),
					logging.Error(err),
					logging.String("drive_id", drive.ID),
				)

				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("command completed",
				logging.String("command", commandName),
				logging.String("drive_id", drive.ID),
			)

			return nil
		},
	}
}
