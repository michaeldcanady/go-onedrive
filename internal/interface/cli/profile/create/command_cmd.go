// Package create provides the command-line interface for initializing new OneDrive profiles.
package create

import (
	"context"
	"strings"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
)

// CreateCmd handles the execution logic for the 'profile create' command.
type CreateCmd struct {
	util.BaseCommand
}

// NewCreateCmd creates a new CreateCmd instance with the provided dependency container.
func NewCreateCmd(container didomain.Container) *CreateCmd {
	return &CreateCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the profile create command. It initializes a new profile with
// the given name, optionally overwriting an existing one and setting it as
// the active profile.
// It uses specific domain services to decouple from the full container.
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

	c.Log.Info("starting profile creation", domainlogger.String("name", name))

	profileSvc := c.Container.Profile()
	if profileSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "profile service is nil")
	}

	c.Log.Debug("checking if profile exists", domainlogger.String("name", name))
	exists, err := profileSvc.Exists(ctx, name)
	if err != nil {
		c.Log.Error("failed to check profile existence", domainlogger.Error(err))
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

		if err := profileSvc.Delete(ctx, name); err != nil {
			c.Log.Error("failed to delete existing profile", domainlogger.Error(err))
			return util.NewCommandErrorWithNameWithError(c.Name, err)
		}
	}

	c.Log.Info("creating profile", domainlogger.String("name", name))

	p, err := profileSvc.Create(ctx, name)
	if err != nil {
		c.Log.Error("failed to create profile", domainlogger.Error(err))
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	if opts.SetCurrent {
		c.Log.Info("setting new profile as current", domainlogger.String("name", name))

		stateSvc := c.Container.State()
		if stateSvc == nil {
			return util.NewCommandErrorWithNameWithMessage(c.Name, "state service is nil")
		}

		if err := stateSvc.Set(domainstate.KeyProfile, p.Name, domainstate.ScopeGlobal); err != nil {
			c.Log.Error("failed to set current profile", domainlogger.Error(err))
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
