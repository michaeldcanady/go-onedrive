package delete

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
)

// Command executes the profile delete operation.
type Command struct {
	profile profile.Service
	log     logger.Logger
}

// NewCommand initializes a new instance of the profile delete Command.
func NewCommand(p profile.Service, l logger.Logger) *Command {
	return &Command{
		profile: p,
		log:     l,
	}
}

// Validate prepares and validates the options for the profile delete operation.
func (c *Command) Validate(ctx context.Context, opts *Options) error {
	return opts.Validate()
}

// Execute deletes a profile.
func (c *Command) Execute(ctx context.Context, opts Options) error {
	log := c.log.WithContext(ctx)

	log.Info("deleting profile", logger.String("name", opts.Name))
	if err := c.profile.Delete(ctx, opts.Name); err != nil {
		log.Error("failed to delete profile", logger.String("name", opts.Name), logger.Error(err))
		return fmt.Errorf("failed to delete profile %s: %w", opts.Name, err)
	}

	log.Info("profile deleted successfully", logger.String("name", opts.Name))
	fmt.Fprintf(opts.Stdout, "Profile '%s' deleted successfully.\n", opts.Name)
	return nil
}

// Finalize performs any necessary cleanup after the profile delete operation.
func (c *Command) Finalize(ctx context.Context, opts Options) error {
	return nil
}
