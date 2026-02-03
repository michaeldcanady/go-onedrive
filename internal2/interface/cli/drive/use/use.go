package use

import (
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
			logger, err := util.EnsureLogger(container, loggerID)
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			id := strings.ToLower(strings.TrimSpace(args[0]))
			if id == "" {
				logger.Warn("id is empty", logging.String("command", commandName))
				return util.NewCommandErrorWithNameWithMessage(commandName, "id is empty")
			}

			drive, err := container.Drive().ResolveDrive(cmd.Context(), id)
			if err != nil {
				logger.Warn("failed to resolve drive", logging.Error(err), logging.String("drive_id", id))
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			if err := container.State().SetCurrentDrive(drive.ID); err != nil {
				logger.Warn("failed to set current drive", logging.Error(err), logging.String("drive_id", id))
				return util.NewCommandErrorWithNameWithMessage(commandName, "failed to set current drive")
			}

			cmd.Printf("Active drive set to %q\n", drive.ID)
			return nil
		},
	}
}
