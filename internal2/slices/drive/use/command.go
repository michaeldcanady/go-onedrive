package use

import (
	"context"
	"fmt"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal2/core/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/core/state"
	"github.com/michaeldcanady/go-onedrive/internal2/di"
	"github.com/spf13/cobra"
)

// Options defines the configuration for the drive use operation.
type Options struct {
	DriveRef string
	Stdout   io.Writer
}

// Handler executes the drive use operation.
type Handler struct {
	drive drive.Service
	state state.Service
	log   logger.Logger
}

// NewHandler initializes a new instance of the drive use Handler.
func NewHandler(drive drive.Service, state state.Service, l logger.Logger) *Handler {
	return &Handler{
		drive: drive,
		state: state,
		log:   l,
	}
}

// Handle sets the specified drive as the active drive for the session.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("switching active drive", logger.String("ref", opts.DriveRef))

	// Resolve the reference (could be alias, ID, or name)
	// For now, our drive.Service.ResolveDrive handles ID and name.
	// We'll also check state for aliases.
	driveID, err := h.state.GetDriveAlias(opts.DriveRef)
	if err != nil {
		// Not an alias, try resolving by name/id
		d, err := h.drive.ResolveDrive(ctx, opts.DriveRef)
		if err != nil {
			return fmt.Errorf("failed to resolve drive '%s': %w", opts.DriveRef, err)
		}
		driveID = d.ID
	}

	if err := h.state.Set(state.KeyDrive, driveID, state.ScopeGlobal); err != nil {
		return fmt.Errorf("failed to set active drive: %w", err)
	}

	fmt.Fprintf(opts.Stdout, "active drive set to '%s'\n", driveID)
	return nil
}

// CreateUseCmd constructs and returns the cobra.Command for the 'drive use' operation.
func CreateUseCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "use <drive-ref>",
		Short: "Set the active OneDrive drive",
		Long:  "Specify a drive ID, name, or alias to be used for subsequent commands.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				DriveRef: args[0],
				Stdout:   cmd.OutOrStdout(),
			}
			log, _ := container.Logger().CreateLogger("drive-use")
			return NewHandler(container.Drive(), container.State(), log).Handle(cmd.Context(), opts)
		},
	}
}
