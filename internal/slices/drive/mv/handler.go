package mv

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs/shared"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
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
	h.log.Info("moving item",
		logger.String("source", opts.Source),
		logger.String("destination", opts.Destination),
	)

	if err := h.manager.Move(ctx, opts.Source, opts.Destination); err != nil {
		return fmt.Errorf("failed to move %s to %s: %w", opts.Source, opts.Destination, err)
	}

	fmt.Fprintf(opts.Stdout, "Moved %s to %s\n", opts.Source, opts.Destination)
	return nil
}
