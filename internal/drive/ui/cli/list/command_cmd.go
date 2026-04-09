package list

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/drive"
	"github.com/michaeldcanady/go-onedrive/internal/drive/alias"
	"github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/fs/formatting"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Handler executes the drive list operation.
type Handler struct {
	drive drive.Service
	alias alias.Service
	log   logger.Logger
}

// NewHandler initializes a new instance of the drive list Handler.
func NewHandler(
	drive drive.Service,
	alias alias.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("drive-list")
	return &Handler{
		drive: drive,
		alias: alias,
		log:   cliLog,
	}
}

// Handle retrieves all available drives and their aliases, then displays them in a table.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx)

	log.Info("listing available drives")

	drives, err := h.drive.ListDrives(ctx)
	if err != nil {
		log.Error(err.Error(), errors.LogFields(err)...)
		return err
	}

	activeDrive, _ := h.drive.GetActive(ctx)

	columns := []formatting.Column{
		formatting.NewColumn(" ", func(item any) string {
			if item.(drive.Drive).ID == activeDrive.ID {
				return "*"
			}
			return ""
		}),
		formatting.NewColumn("Alias", func(item any) string {
			alias, err := h.alias.GetAliasByDriveID(item.(drive.Drive).Name)
			if err != nil {
				log.Debug("failed to get alias for drive", logger.String("drive_id", item.(drive.Drive).ID), logger.Error(err))
			}

			return alias
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
