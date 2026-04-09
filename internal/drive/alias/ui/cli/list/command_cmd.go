package list

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/drive/alias"
	"github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/fs/formatting"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Handler executes the drive alias list operation.
type Handler struct {
	alias alias.Service
	log   logger.Logger
}

// NewHandler initializes a new instance of the drive alias list Handler.
func NewHandler(
	alias alias.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("alias-list")
	return &Handler{
		alias: alias,
		log:   cliLog,
	}
}

type aliasEntry struct {
	Alias   string
	DriveID string
}

// Handle retrieves and displays all registered drive aliases.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx)

	aliases, err := h.alias.ListAliases()
	if err != nil {
		log.Error(err.Error(), errors.LogFields(err)...)
		return err
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
