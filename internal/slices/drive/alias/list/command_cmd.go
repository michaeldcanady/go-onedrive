package list

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs/formatting"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/core/state"
)

// Handler executes the drive alias list operation.
type Handler struct {
	state state.Service
	log   logger.Logger
}

// NewHandler initializes a new instance of the drive alias list Handler.
func NewHandler(state state.Service, l logger.Logger) *Handler {
	return &Handler{
		state: state,
		log:   l,
	}
}

type aliasEntry struct {
	Alias   string
	DriveID string
}

// Handle retrieves and displays all registered drive aliases.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	aliases, err := h.state.ListDriveAliases()
	if err != nil {
		return fmt.Errorf("failed to list drive aliases: %w", err)
	}

	var entries []any
	for k, v := range aliases {
		entries = append(entries, aliasEntry{Alias: k, DriveID: v})
	}

	columns := []formatting.Column{
		formatting.NewColumn("Alias", func(item any) string { return item.(aliasEntry).Alias }),
		formatting.NewColumn("DriveID", func(item any) string { return item.(aliasEntry).DriveID }),
	}

	formatter := formatting.NewTableFormatter(columns...)
	return formatter.Format(opts.Stdout, entries)
}
