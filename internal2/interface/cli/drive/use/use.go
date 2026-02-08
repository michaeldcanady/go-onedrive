package use

import (
	"context"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/spf13/cobra"
)

const (
	commandName = "use"
	loggerID    = "cli"
)

func CreateUseCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "use <id>",
		Short: "Set current drive",
		Args:  cobra.ExactArgs(1),

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
				logger.Warn("drive id is empty",
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

			logger.Info("setting current drive",
				logging.String("event", "set_current_drive"),
				logging.String("drive_id", drive.ID),
			)

			if err := container.State().SetCurrentDrive(drive.ID); err != nil {
				logger.Warn("failed to set current drive",
					logging.String("event", "set_current_drive"),
					logging.Error(err),
					logging.String("drive_id", drive.ID),
				)
				return util.NewCommandErrorWithNameWithMessage(commandName, "failed to set current drive")
			}

			logger.Info("current drive updated",
				logging.String("event", "drive_selected"),
				logging.String("drive_id", drive.ID),
			)

			cmd.Printf("Active drive set to %q\n", drive.ID)

			logger.Info("command completed",
				logging.String("command", commandName),
				logging.String("drive_id", drive.ID),
			)

			return nil
		},
	}
}
