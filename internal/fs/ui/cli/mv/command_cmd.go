package mv

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Handler executes the drive mv operation.
type Handler struct {
	manager shared.Service
	log     logger.Logger
}

// NewHandler initializes a new instance of the drive mv Handler.
func NewHandler(
	m shared.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("drive-mv")
	return &Handler{
		manager: m,
		log:     cliLog,
	}
}

// Handle moves an item from the source to the destination.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx).With(
		logger.String("source", opts.Source),
		logger.String("destination", opts.Destination),
	)

	log.Info("starting move operation")

	log.Debug("delegating to filesystem manager for move")
	if err := h.manager.Move(ctx, opts.Source, opts.Destination); err != nil {
		wrapped := cli.WrapError(err, opts.Source)
		h.log.Error(wrapped.Error(), errors.LogFields(wrapped)...)
		return wrapped
	}

	log.Info("move completed successfully")
	fmt.Fprintf(opts.Stdout, "Moved %s to %s\n", opts.Source, opts.Destination)
	return nil
}
