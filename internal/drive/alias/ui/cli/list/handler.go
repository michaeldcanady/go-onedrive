package list

import (
	"context"
	"fmt"
	"sort"

	"github.com/michaeldcanady/go-onedrive/internal/drive/alias"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Command executes the drive alias list operation.
type Command struct {
	alias alias.Service
	log   logger.Logger
}

// NewCommand initializes a new instance of the drive alias list Command.
func NewCommand(a alias.Service, l logger.Logger) *Command {
	return &Command{
		alias: a,
		log:   l,
	}
}

// Validate prepares and validates the options for the drive alias list operation.
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	return opts.Validate()
}

// Execute retrieves and displays all configured drive aliases.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx)

	log.Debug("fetching all aliases")
	aliases, err := c.alias.ListAliases()
	if err != nil {
		log.Error("failed to list aliases", logger.Error(err))
		return fmt.Errorf("failed to list aliases: %w", err)
	}

	if len(aliases) == 0 {
		log.Info("no aliases found")
		fmt.Fprintln(opts.Stdout, "No drive aliases configured.")
		return nil
	}

	log.Info("aliases retrieved successfully", logger.Int("count", len(aliases)))

	// Sort aliases by name for consistent output
	names := make([]string, 0, len(aliases))
	for _, name := range aliases {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Fprintln(opts.Stdout, "Configured drive aliases:")
	for _, name := range names {
		// ListAliases returns map[driveID]aliasName
		// We need to find the driveID for this name
		for id, n := range aliases {
			if n == name {
				fmt.Fprintf(opts.Stdout, "  %s -> %s\n", name, id)
				break
			}
		}
	}

	return nil
}

// Finalize performs any necessary cleanup after the drive alias list operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
