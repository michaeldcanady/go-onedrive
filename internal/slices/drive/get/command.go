package get

import (
	"context"
	"fmt"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/core/drive"
	"github.com/michaeldcanady/go-onedrive/internal/core/fs/formatting"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/core/state"
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// Options defines the configuration for the drive get operation.
type Options struct {
	DriveRef string
	Stdout   io.Writer
}

// Handler executes the drive get operation.
type Handler struct {
	drive drive.Service
	state state.Service
	log   logger.Logger
}

// NewHandler initializes a new instance of the drive get Handler.
func NewHandler(drive drive.Service, state state.Service, l logger.Logger) *Handler {
	return &Handler{
		drive: drive,
		state: state,
		log:   l,
	}
}

// Handle retrieves and displays details for a specific drive.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("fetching drive details", logger.String("ref", opts.DriveRef))

	var driveID string
	if opts.DriveRef == "" {
		// Use active drive
		var err error
		driveID, err = h.state.Get(state.KeyDrive)
		if err != nil {
			return fmt.Errorf("no active drive set and no drive reference provided: %w", err)
		}
	} else {
		// Resolve the reference
		var err error
		driveID, err = h.state.GetDriveAlias(opts.DriveRef)
		if err != nil {
			// Not an alias, assume it's ID or name
			driveID = opts.DriveRef
		}
	}

	d, err := h.drive.ResolveDrive(ctx, driveID)
	if err != nil {
		return fmt.Errorf("failed to get drive details: %w", err)
	}

	columns := []formatting.Column{
		formatting.NewColumn("ID", func(item any) string { return item.(drive.Drive).ID }),
		formatting.NewColumn("Name", func(item any) string { return item.(drive.Drive).Name }),
		formatting.NewColumn("Type", func(item any) string { return item.(drive.Drive).Type.String() }),
		formatting.NewColumn("Owner", func(item any) string { return item.(drive.Drive).Owner }),
		formatting.NewColumn("ReadOnly", func(item any) string { return fmt.Sprintf("%v", item.(drive.Drive).ReadOnly) }),
	}

	formatter := formatting.NewTableFormatter(columns...)
	return formatter.Format(opts.Stdout, []any{d})
}

// CreateGetCmd constructs and returns the cobra.Command for the 'drive get' operation.
func CreateGetCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "get [drive-ref]",
		Short: "Display details for a specific drive",
		Long:  "Provide a drive ID, name, or alias. If omitted, the currently active drive is shown.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			driveRef := ""
			if len(args) > 0 {
				driveRef = args[0]
			}
			opts := Options{
				DriveRef: driveRef,
				Stdout:   cmd.OutOrStdout(),
			}
			log, _ := container.Logger().CreateLogger("drive-get")
			return NewHandler(container.Drive(), container.State(), log).Handle(cmd.Context(), opts)
		},
	}
}
