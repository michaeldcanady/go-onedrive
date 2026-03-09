package list

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/formatting"
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

			stateSvc := container.State()
			activeDriveID, _ := stateSvc.GetCurrentDrive()
			aliases, _ := stateSvc.ListDriveAliases()

			// Prepare alias lookup
			aliasMap := make(map[string]string)
			for alias, driveID := range aliases {
				aliasMap[driveID] = alias
			}

			columns := []formatting.Column[*drive.Drive]{
				formatting.NewColumn(" ", func(item *drive.Drive) string {
					if item.ID == activeDriveID {
						return "*"
					}
					return ""
				}),
				formatting.NewColumn("Alias", func(item *drive.Drive) string {
					return aliasMap[item.ID]
				}),
				formatting.NewColumn("ID", func(item *drive.Drive) string { return item.ID }),
				formatting.NewColumn("Name", func(item *drive.Drive) string { return item.Name }),
				formatting.NewColumn("Type", func(item *drive.Drive) string { return string(item.Type) }),
			}

			formatter := formatting.NewTableFormatter(columns...).WithTruncate(true)

			if err := formatter.Format(cmd.OutOrStdout(), drives); err != nil {
				logger.Warn("failed to format output", logging.Error(err))
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			return nil
		},
	}
}
