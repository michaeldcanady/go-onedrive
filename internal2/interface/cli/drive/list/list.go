package list

import (
	"context"
	"reflect"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	domaindrive "github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			_, err := util.EnsureLogger(container, loggerID)
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			tableFormatter := NewTableFormatter(driveIDColumn, driveNameColumn, driveOwnerColumn, driveReadOnlyColumn, driveTypeColumn)
			if err := container.Format().RegisterWithType("table", reflect.TypeFor[[]*domaindrive.Drive](), tableFormatter); err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			return nil
		},
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

			if err := container.Format().Format(cmd.OutOrStdout(), "table", drives); err != nil {
				logger.Warn("failed to format output", logging.Error(err))
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			return nil
		},
	}
}
