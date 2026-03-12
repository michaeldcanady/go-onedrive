package remove

import (
	"context"
	"fmt"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal2/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/core/state"
	"github.com/michaeldcanady/go-onedrive/internal2/di"
	"github.com/spf13/cobra"
)

// Options defines the configuration for the drive alias remove operation.
type Options struct {
	Alias  string
	Stdout io.Writer
}

// Handler executes the drive alias remove operation.
type Handler struct {
	state state.Service
	log   logger.Logger
}

// NewHandler initializes a new instance of the drive alias remove Handler.
func NewHandler(state state.Service, l logger.Logger) *Handler {
	return &Handler{
		state: state,
		log:   l,
	}
}

// Handle deletes a drive alias.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("removing drive alias", logger.String("alias", opts.Alias))

	if err := h.state.RemoveDriveAlias(opts.Alias); err != nil {
		return fmt.Errorf("failed to remove drive alias: %w", err)
	}

	fmt.Fprintf(opts.Stdout, "alias '%s' removed\n", opts.Alias)
	return nil
}

// CreateRemoveCmd constructs and returns the cobra.Command for the 'drive alias remove' operation.
func CreateRemoveCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <alias>",
		Short: "Remove a drive alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				Alias:  args[0],
				Stdout: cmd.OutOrStdout(),
			}
			log, _ := container.Logger().CreateLogger("alias-remove")
			return NewHandler(container.State(), log).Handle(cmd.Context(), opts)
		},
	}
}
