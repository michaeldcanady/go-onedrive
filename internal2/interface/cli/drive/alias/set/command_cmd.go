package set

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type Command struct {
	container di.Container
	logger    infralogging.Logger
}

func NewCmd(container di.Container) *Command {
	return &Command{
		container: container,
	}
}

func (c *Command) Run(ctx context.Context, opts Options) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if c.logger == nil {
		logger, err := util.EnsureLogger(c.container, loggerID)
		if err != nil {
			return util.NewCommandErrorWithNameWithError(commandName, err)
		}
		c.logger = logger
	}

	c.logger.Info("setting drive alias",
		infralogging.String("alias", opts.Alias),
		infralogging.String("drive_id", opts.DriveID),
	)

	if err := c.container.State().SetDriveAlias(opts.Alias, opts.DriveID); err != nil {
		c.logger.Error("failed to set drive alias",
			infralogging.Error(err),
		)
		return util.NewCommandError(commandName, "failed to set drive alias", err)
	}

	fmt.Fprintf(opts.Stdout, "Alias %q set to drive %q\n", opts.Alias, opts.DriveID)

	return nil
}
