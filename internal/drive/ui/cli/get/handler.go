package get

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/drive"
	"github.com/michaeldcanady/go-onedrive/internal/drive/alias"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Command executes the drive get operation.
type Command struct {
	drive drive.Service
	alias alias.Service
	log   logger.Logger
}

// NewCommand initializes a new instance of the drive get Command.
func NewCommand(d drive.Service, a alias.Service, l logger.Logger) *Command {
	return &Command{
		drive: d,
		alias: a,
		log:   l,
	}
}

// Validate prepares and validates the options for the drive get operation.
func (c *Command) Validate(ctx *CommandContext) error {
	return ctx.Options.Validate()
}

// Execute retrieves and displays the active drive.
func (c *Command) Execute(ctx *CommandContext) error {
	log := c.log.WithContext(ctx.Ctx)

	log.Debug("fetching active drive", logger.String("identity", ctx.Options.IdentityID))
	d, err := c.drive.GetActive(ctx.Ctx, ctx.Options.IdentityID)
	if err != nil {
		log.Error("failed to get active drive", logger.Error(err))
		return fmt.Errorf("failed to get active drive: %w", err)
	}

	log.Debug("fetching aliases for active drive")
	aliases, err := c.alias.ListAliases()
	if err != nil {
		log.Warn("failed to list aliases", logger.Error(err))
	}

	var activeAliases []string
	for id, name := range aliases {
		if id == d.ID {
			activeAliases = append(activeAliases, name)
		}
	}

	log.Info("active drive retrieved successfully", logger.String("id", d.ID))
	fmt.Fprintf(ctx.Options.Stdout, "Active drive: %s (%s)\n", d.Name, d.ID)
	if len(activeAliases) > 0 {
		fmt.Fprintf(ctx.Options.Stdout, "Aliases: %s\n", fmt.Sprint(activeAliases))
	}

	return nil
}

// Finalize performs any necessary cleanup after the drive get operation.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
