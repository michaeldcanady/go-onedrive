package current

import (
	"context"
	"fmt"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
)

type CurrentCmd struct {
	util.BaseCommand
}

func NewCurrentCmd(container didomain.Container) *CurrentCmd {
	return &CurrentCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *CurrentCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("retrieving current profile")

	name, err := c.Container.State().GetCurrentProfile()
	if err != nil {
		c.RenderError(opts.Stderr, err)
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	fmt.Fprintf(opts.Stdout, "%s\n", name)

	c.Log.Info("profile current completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
