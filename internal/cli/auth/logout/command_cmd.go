package logout

import (
	"context"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
)

type LogoutCmd struct {
	util.BaseCommand
}

func NewLogoutCmd(container didomain.Container) *LogoutCmd {
	return &LogoutCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *LogoutCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("starting logout flow")

	profileName, err := c.Container.State().GetCurrentProfile()
	if err != nil {
		c.Log.Error("failed to get current profile", domainlogger.Error(err))
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	c.Log.Info("resolved current profile",
		domainlogger.String("profile", profileName),
		domainlogger.Bool("force", opts.Force),
	)

	authService := c.Container.Auth()

	c.Log.Info("attempting logout",
		domainlogger.String("profile", profileName),
	)

	err = authService.Logout(ctx, profileName, opts.Force)
	if err != nil {
		c.Log.Error("logout failed",
			domainlogger.String("profile", profileName),
			domainlogger.Error(err),
		)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	c.Log.Info("logout successful",
		domainlogger.String("profile", profileName),
	)

	c.RenderSuccess(opts.Stdout, "Logged out of profile %q\n", profileName)

	c.Log.Info("logout flow completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
