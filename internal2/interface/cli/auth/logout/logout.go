package logout

import (
	"context"

	applogging "github.com/michaeldcanady/go-onedrive/internal2/app/common/logging"
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
			}

			logger, err := ensureLogger(container)
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("starting logout flow")

			// Determine which profile to log out from
			profileName, err := container.State().GetCurrentProfile()
			if err != nil {
				logger.Error("failed to get current profile", infralogging.Error(err))
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("resolved current profile",
				infralogging.String("profile", profileName),
				infralogging.Bool("force", force),
			)

			authService := container.Auth()

			// Attempt logout
			logger.Info("attempting logout",
				infralogging.String("profile", profileName),
			)

			err = authService.Logout(ctx, profileName, force)
			if err != nil {
				logger.Error("logout failed",
					infralogging.String("profile", profileName),
					infralogging.Error(err),
				)
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("logout successful",
				infralogging.String("profile", profileName),
			)

			cmd.Printf("Logged out of profile %q\n", profileName)
			logger.Info("logout flow completed successfully")

			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, forceLongFlag, forceShortFlag, false, forceUsage)

	return cmd
}

func ensureLogger(c di.Container) (infralogging.Logger, error) {
	logger, err := c.Logger().GetLogger(loggerID)
	if err == applogging.ErrUnknownLogger {
		return c.Logger().CreateLogger(loggerID)
	}
	return logger, err
}
