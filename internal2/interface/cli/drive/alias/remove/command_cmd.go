package remove

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

	c.logger.Info("removing drive alias",
		infralogging.String("alias", opts.Alias),
	)

	if err := c.container.State().RemoveDriveAlias(opts.Alias); err != nil {
		c.logger.Error("failed to remove drive alias",
			infralogging.Error(err),
		)
		return util.NewCommandError(commandName, "failed to remove drive alias", err)
	}

	fmt.Fprintf(opts.Stdout, "Alias %q removed\n", opts.Alias)

	return nil
}
