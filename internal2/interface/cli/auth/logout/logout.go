package logout

import (
	"context"
	"fmt"
	"time"

	logger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
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

type LogoutCmd struct {
	util.BaseCommand
}

func NewLogoutCmd(container di.Container) *LogoutCmd {
	return &LogoutCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *LogoutCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("starting logout flow")

	profileName, err := c.Container.State().GetCurrentProfile()
	if err != nil {
		c.Log.Error("failed to get current profile", logger.Error(err))
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	c.Log.Info("resolved current profile",
		logger.String("profile", profileName),
		logger.Bool("force", opts.Force),
	)

	authService := c.Container.Auth()

	c.Log.Info("attempting logout",
		logger.String("profile", profileName),
	)

	err = authService.Logout(ctx, profileName, opts.Force)
	if err != nil {
		c.Log.Error("logout failed",
			logger.String("profile", profileName),
			logger.Error(err),
		)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	c.Log.Info("logout successful",
		logger.String("profile", profileName),
	)

	fmt.Fprintf(opts.Stdout, "Logged out of profile %q\n", profileName)

	c.Log.Info("logout flow completed successfully",
		logger.Duration("duration", time.Since(start)),
	)
	return nil
}

func CreateLogoutCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out of the current OneDrive profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			return NewLogoutCmd(container).Run(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Force, forceLongFlag, forceShortFlag, false, forceUsage)

	return cmd
}
