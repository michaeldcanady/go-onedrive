package use

import (
	"context"
	"strings"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
)

type UseCmd struct {
	util.BaseCommand
}

func NewUseCmd(container didomain.Container) *UseCmd {
	return &UseCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *UseCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	name := strings.TrimSpace(opts.Name)
	if name == "" {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "name is empty")
	}

	name = strings.ToLower(name)

	c.Log.Info("setting current profile", domainlogger.String("profile", name))

	// Validate profile exists
	p, err := c.Container.Profile().Get(ctx, name)
	if err != nil {
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	// Persist as current profile
	if err := c.Container.State().SetCurrentProfile(p.Name, domainstate.ScopeGlobal); err != nil {
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	c.RenderSuccess(opts.Stdout, "active profile set to %q", p.Name)

	c.Log.Info("profile use completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
