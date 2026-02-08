package logout

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/spf13/cobra"
)

const (
	commandName = "logout"
	loggerID    = "cli"

	forceLongFlag  = "force"
	forceShortFlag = "f"
	forceUsage     = "Force logout even if no active session is detected"
)

func CreateLogoutCmd(container di.Container) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out of the current OneDrive profile",

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
				infralogging.String("command", commandName),
				infralogging.Bool("force", force),
			)

			logger.Debug("resolving current profile",
				infralogging.String("event", "resolve_profile"),
			)

			profileName, err := container.State().GetCurrentProfile()
			if err != nil {
				logger.Error("failed to get current profile",
					infralogging.String("event", "resolve_profile"),
					infralogging.Error(err),
				)
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("resolved current profile",
				infralogging.String("event", "profile_resolved"),
				infralogging.String("profile", profileName),
				infralogging.Bool("force", force),
			)

			authService := container.Auth()

			logger.Info("attempting logout",
				infralogging.String("event", "logout_attempt"),
				infralogging.String("profile", profileName),
				infralogging.Bool("force", force),
			)

			err = authService.Logout(ctx, profileName, force)
			if err != nil {
				logger.Error("logout failed",
					infralogging.String("event", "logout_failed"),
					infralogging.String("profile", profileName),
					infralogging.Bool("force", force),
					infralogging.Error(err),
				)
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("logout successful",
				infralogging.String("event", "logout_success"),
				infralogging.String("profile", profileName),
				infralogging.Bool("force", force),
			)

			cmd.Printf("Logged out of profile %q\n", profileName)

			logger.Info("command completed",
				infralogging.String("command", commandName),
				infralogging.String("profile", profileName),
				infralogging.Bool("force", force),
			)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, forceLongFlag, forceShortFlag, false, forceUsage)

	return cmd
}
