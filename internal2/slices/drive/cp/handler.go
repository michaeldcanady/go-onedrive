package cp

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/shared"
	"github.com/michaeldcanady/go-onedrive/internal2/core/logger"
)

// Handler executes the drive cp operation.
type Handler struct {
	manager shared.Service
	log     logger.Logger
}

// NewHandler initializes a new instance of the drive cp Handler.
func NewHandler(m shared.Service, l logger.Logger) *Handler {
	return &Handler{
		manager: m,
		log:     l,
	}
}

// Handle copies an item from the source to the destination.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("copying item",
		logger.String("source", opts.Source),
		logger.String("destination", opts.Destination),
		logger.Bool("recursive", opts.Recursive),
	)

	cpOpts := shared.CopyOptions{
		Recursive: opts.Recursive,
		Overwrite: true, // Default to overwrite for now, can be a flag later.
	}

	if err := h.manager.Copy(ctx, opts.Source, opts.Destination, cpOpts); err != nil {
		return fmt.Errorf("failed to copy %s to %s: %w", opts.Source, opts.Destination, err)
	}

	fmt.Fprintf(opts.Stdout, "Copied %s to %s\n", opts.Source, opts.Destination)
	return nil
}
