package cp

import (
	"context"
	"fmt"

	shared "github.com/michaeldcanady/go-onedrive/internal/feature/fs"
	"github.com/michaeldcanady/go-onedrive/internal/feature/logger"
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
	log := h.log.WithContext(ctx).With(
		logger.String("source", opts.Source),
		logger.String("destination", opts.Destination),
		logger.Bool("recursive", opts.Recursive),
	)

	log.Info("starting copy operation")

	cpOpts := shared.CopyOptions{
		Recursive: opts.Recursive,
		Overwrite: true, // Default to overwrite for now, can be a flag later.
	}

	log.Debug("delegating to filesystem manager for copy")
	if err := h.manager.Copy(ctx, opts.Source, opts.Destination, cpOpts); err != nil {
		log.Error("copy failed", logger.Error(err))
		return fmt.Errorf("failed to copy %s to %s: %w", opts.Source, opts.Destination, err)
	}

	log.Info("copy completed successfully")
	fmt.Fprintf(opts.Stdout, "Copied %s to %s\n", opts.Source, opts.Destination)
	return nil
}
