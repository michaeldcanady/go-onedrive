package list

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/feature/drive"
	"github.com/michaeldcanady/go-onedrive/internal/feature/fs/formatting"
	"github.com/michaeldcanady/go-onedrive/internal/feature/logger"
	"github.com/michaeldcanady/go-onedrive/internal/feature/state"
)

// Handler executes the drive list operation.
type Handler struct {
	drive drive.Service
	state state.Service
	log   logger.Logger
}

// NewHandler initializes a new instance of the drive list Handler.
func NewHandler(drive drive.Service, state state.Service, l logger.Logger) *Handler {
	return &Handler{
		drive: drive,
		state: state,
		log:   l,
	}
}

// Handle retrieves all available drives and their aliases, then displays them in a table.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("listing available drives")

	drives, err := h.drive.ListDrives(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve drives: %w", err)
	}

	activeDriveID, _ := h.state.Get(state.KeyDrive)
	aliases, _ := h.state.ListDriveAliases()

	// Prepare alias lookup
	aliasMap := make(map[string]string)
	for alias, driveID := range aliases {
		aliasMap[driveID] = alias
	}

	columns := []formatting.Column{
		formatting.NewColumn(" ", func(item any) string {
			if item.(drive.Drive).ID == activeDriveID {
				return "*"
			}
			return ""
		}),
		formatting.NewColumn("Alias", func(item any) string {
			return aliasMap[item.(drive.Drive).ID]
		}),
		formatting.NewColumn("ID", func(item any) string { return item.(drive.Drive).ID }),
		formatting.NewColumn("Name", func(item any) string { return item.(drive.Drive).Name }),
		formatting.NewColumn("Type", func(item any) string { return item.(drive.Drive).Type.String() }),
	}

	anyDrives := make([]any, len(drives))
	for i, d := range drives {
		anyDrives[i] = d
	}

	formatter := formatting.NewTableFormatter(columns...)
	return formatter.Format(opts.Stdout, anyDrives)
}
