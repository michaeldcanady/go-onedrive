package list

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			logger, err := util.EnsureLogger(container, loggerID)
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			drives, err := container.Drive().ListDrives(ctx)
			if err != nil {
				logger.Warn("failed to retrieve drives", logging.Error(err))
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			formatter := NewTableFormatter(driveIDColumn, driveNameColumn, driveOwnerColumn, driveReadOnlyColumn, driveTypeColumn)
			if err := formatter.Format(cmd.OutOrStdout(), drives); err != nil {
				logger.Warn("failed to format output", logging.Error(err))
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			return nil
		},
	}
}
