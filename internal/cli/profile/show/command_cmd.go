package show

import (
	"context"
	"fmt"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
)

type ShowCmd struct {
	util.BaseCommand
}

func NewShowCmd(container didomain.Container) *ShowCmd {
	return &ShowCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *ShowCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("showing profile details", domainlogger.String("profile", opts.Name))

	p, err := c.Container.Profile().Get(ctx, opts.Name)
	if err != nil {
		c.RenderError(opts.Stderr, err)
		return err
	}

	fmt.Fprintf(opts.Stdout, "Name: %s\n", p.Name)
	fmt.Fprintf(opts.Stdout, "Path: %s\n", p.Path)
	if p.ConfigurationPath != "" {
		fmt.Fprintf(opts.Stdout, "Config: %s\n", p.ConfigurationPath)
	}

	c.Log.Info("profile show completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
