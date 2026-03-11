package create

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
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
		return err
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

		if err := c.Container.State().SetCurrentProfile(p.Name); err != nil {
			c.Log.Error("failed to set current profile", domainlogger.String("error", err.Error()))
			return util.NewCommandErrorWithNameWithError(c.Name, err)
		}
	}

	c.Log.Info("profile created successfully",
		domainlogger.String("name", p.Name),
		domainlogger.String("path", p.Path),
	)

	fmt.Fprintf(opts.Stdout, "Created profile %q at %s\n", p.Name, p.Path)

	c.Log.Info("profile create completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
