package current

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/spf13/cobra"
)

const (
	commandName = "current"
	loggerID    = "cli"
)

func CreateCurrentCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show the active profile",

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

			logger.Debug("retrieving current profile",
				logging.String("event", "get_current_profile"),
			)

			name, err := container.State().GetCurrentProfile()
			if err != nil {
				logger.Warn("failed to retrieve current profile",
					logging.String("event", "get_current_profile"),
					logging.Error(err),
				)
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("current profile retrieved",
				logging.String("event", "profile_resolved"),
				logging.String("profile", name),
			)

			cmd.Printf("%s\n", name)

			logger.Info("command completed",
				logging.String("command", commandName),
				logging.String("profile", name),
			)

			return nil
		},
	}
}
