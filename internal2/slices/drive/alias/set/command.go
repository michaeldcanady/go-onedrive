package set

import (
	"context"
	"fmt"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal2/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/core/state"
	"github.com/michaeldcanady/go-onedrive/internal2/di"
	"github.com/spf13/cobra"
)

// Options defines the configuration for the drive alias set operation.
type Options struct {
	Alias   string
	DriveID string
	Stdout  io.Writer
}

// Handler executes the drive alias set operation.
type Handler struct {
	state state.Service
	log   logger.Logger
}

// NewHandler initializes a new instance of the drive alias set Handler.
func NewHandler(state state.Service, l logger.Logger) *Handler {
	return &Handler{
		state: state,
		log:   l,
	}
}

// Handle creates or updates a drive alias.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("setting drive alias", logger.String("alias", opts.Alias), logger.String("driveID", opts.DriveID))

	if err := h.state.SetDriveAlias(opts.Alias, opts.DriveID); err != nil {
		return fmt.Errorf("failed to set drive alias: %w", err)
	}

	fmt.Fprintf(opts.Stdout, "alias '%s' set to drive '%s'\n", opts.Alias, opts.DriveID)
	return nil
}

// CreateSetCmd constructs and returns the cobra.Command for the 'drive alias set' operation.
func CreateSetCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "set <alias> <drive-id>",
		Short: "Create or update a drive alias",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				Alias:   args[0],
				DriveID: args[1],
				Stdout:  cmd.OutOrStdout(),
			}
			log, _ := container.Logger().CreateLogger("alias-set")
			return NewHandler(container.State(), log).Handle(cmd.Context(), opts)
		},
	}
}
