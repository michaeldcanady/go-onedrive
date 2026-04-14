package mv

import (
	"context"
	"fmt"

	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Handler executes the drive mv operation.
type Handler struct {
	manager shared.Service
	log     logger.Logger
}

// NewHandler initializes a new instance of the drive mv Handler.
func NewHandler(m shared.Service, l logger.Logger) *Handler {
	return &Handler{
		manager: m,
		log:     l,
	}
}

// Handle moves an item from the source to the destination.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx).With(
		logger.String("source", opts.SourceURI.String()),
		logger.String("destination", opts.DestinationURI.String()),
	)

	log.Info("starting move operation")

	log.Debug("delegating to filesystem manager for move")
	if err := h.manager.Move(ctx, opts.SourceURI, opts.DestinationURI); err != nil {
		log.Error("move failed", logger.Error(err))
		return fmt.Errorf("failed to move %s to %s: %w", opts.SourceURI, opts.DestinationURI, err)
	}

	log.Info("move completed successfully")
	fmt.Fprintf(opts.Stdout, "Moved %s to %s\n", opts.SourceURI, opts.DestinationURI)
	return nil
}
