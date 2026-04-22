package create

import (
	"errors"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/features/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/profile"
)

// Command executes the profile create operation.
type Command struct {
	profile profile.Service
	log     logger.Logger
}

// NewCommand initializes a new instance of the profile create Command.
func NewCommand(p profile.Service, l logger.Logger) *Command {
	return &Command{
		profile: p,
		log:     l,
	}
}

// Validate prepares and validates the options for the profile create operation.
func (c *Command) Validate(ctx *CommandContext) error {
	ctx.Options.Name = strings.TrimSpace(ctx.Options.Name)
	if ctx.Options.Name == "" {
		return errors.New("profile name is required")
	}

	return nil
}

// Execute creates a new profile.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)

	log.Info("creating profile", logger.String("name", ctx.Options.Name))
	if profile, err := c.profile.Create(ctx.Ctx, ctx.Options.Name); err != nil {
		log.Error("failed to create profile", logger.String("name", ctx.Options.Name), logger.Error(err))
		return fmt.Errorf("failed to create profile %s: %w", ctx.Options.Name, err)
	} else {
		ctx.Profile = profile
	}

	log.Info("profile created successfully", logger.String("name", ctx.Options.Name))
	return nil
}

// Finalize performs any necessary cleanup after the profile create operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	_, _ = fmt.Fprintf(ctx.Options.Stdout, "Profile '%s' created successfully.\n", ctx.Profile.Name)

	return nil
}
