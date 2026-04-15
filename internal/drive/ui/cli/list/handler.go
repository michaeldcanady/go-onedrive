package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/drive"
	"github.com/michaeldcanady/go-onedrive/internal/drive/alias"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Command executes the drive list operation.
type Command struct {
	drive drive.Service
	alias alias.Service
	log   logger.Logger
}

// NewCommand initializes a new instance of the drive list Command.
func NewCommand(d drive.Service, a alias.Service, l logger.Logger) *Command {
	return &Command{
		drive: d,
		alias: a,
		log:   l,
	}
}

// Validate prepares and validates the options for the drive list operation.
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	return opts.Validate()
}

// Execute retrieves and displays all available OneDrive drives.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx)

	log.Debug("fetching all drives")
	drives, err := c.drive.ListDrives(ctx)
	if err != nil {
		log.Error("failed to list drives", logger.Error(err))
		return fmt.Errorf("failed to list drives: %w", err)
	}

	log.Debug("fetching active drive for marking")
	current, _ := c.drive.GetActive(ctx)

	log.Debug("fetching aliases")
	aliases, _ := c.alias.ListAliases()

	log.Info("drives retrieved successfully", logger.Int("count", len(drives)))
	fmt.Fprintln(opts.Stdout, "Available OneDrive drives:")
	for _, d := range drives {
		prefix := "  "
		if d.ID == current.ID {
			prefix = "* "
		}

		var driveAliases []string
		for id, name := range aliases {
			if id == d.ID {
				driveAliases = append(driveAliases, name)
			}
		}

		aliasStr := ""
		if len(driveAliases) > 0 {
			aliasStr = fmt.Sprintf(" [Aliases: %s]", strings.Join(driveAliases, ", "))
		}

		fmt.Fprintf(opts.Stdout, "%s%s (%s)%s\n", prefix, d.Name, d.ID, aliasStr)
	}

	return nil
}

// Finalize performs any necessary cleanup after the drive list operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
