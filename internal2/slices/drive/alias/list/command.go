package list

import (
	"context"
	"fmt"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/formatting"
	"github.com/michaeldcanady/go-onedrive/internal2/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/core/state"
	"github.com/michaeldcanady/go-onedrive/internal2/di"
	"github.com/spf13/cobra"
)

// Options defines the configuration for the drive alias list operation.
type Options struct {
	// Stdout is the writer for standard output.
	Stdout io.Writer
}

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

// CreateListCmd constructs and returns the cobra.Command for the 'drive alias list' operation.
func CreateListCmd(container di.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all drive aliases",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{Stdout: cmd.OutOrStdout()}
			log, _ := container.Logger().CreateLogger("alias-list")
			return NewHandler(container.State(), log).Handle(cmd.Context(), opts)
		},
	}
}
