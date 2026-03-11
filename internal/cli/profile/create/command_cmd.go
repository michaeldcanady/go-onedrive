package create

import (
	"context"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
	"strings"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
)

type CreateCmd struct {
	util.BaseCommand
}

func NewCreateCmd(container didomain.Container) *CreateCmd {
	return &CreateCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *CreateCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	name := strings.ToLower(strings.TrimSpace(opts.Name))
	if name == "" {
		c.Log.Warn("profile name is empty")
		return util.NewCommandErrorWithNameWithMessage(c.Name, "name is empty")
	}

	c.Log.Info("checking if profile exists", domainlogger.String("name", name))

	exists, err := c.Container.Profile().Exists(ctx, name)
	if err != nil {
		c.Log.Error("failed to check profile existence", domainlogger.String("error", err.Error()))
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	// Handle existing profile
	if exists {
		if !opts.Force {
			c.Log.Warn("profile already exists", domainlogger.String("name", name))
			return util.NewCommandErrorWithNameWithMessage(c.Name, "profile already exists")
		}

		c.Log.Warn("profile exists; force enabled, deleting existing profile",
			domainlogger.String("name", name),
		)

		if err := c.Container.Profile().Delete(ctx, name); err != nil {
			c.Log.Error("failed to delete existing profile", domainlogger.String("error", err.Error()))
			return util.NewCommandErrorWithNameWithError(c.Name, err)
		}
	}

	c.Log.Info("creating profile", domainlogger.String("name", name))

	p, err := c.Container.Profile().Create(ctx, name)
	if err != nil {
		c.Log.Error("failed to create profile", domainlogger.String("error", err.Error()))
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	if opts.SetCurrent {
		c.Log.Info("setting new profile as current", domainlogger.String("name", name))

		if err := c.Container.State().Set(domainstate.KeyProfile, p.Name, domainstate.ScopeGlobal); err != nil {
			c.Log.Error("failed to set current profile", domainlogger.String("error", err.Error()))
			return util.NewCommandErrorWithNameWithError(c.Name, err)
		}
	}

	c.Log.Info("profile created successfully",
		domainlogger.String("name", p.Name),
		domainlogger.String("path", p.Path),
	)

	c.RenderSuccess(opts.Stdout, "created profile %q at %s", p.Name, p.Path)

	c.Log.Info("profile create completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
