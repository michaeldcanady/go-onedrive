package download

import (
	"context"
	"fmt"
	"strings"

	shared "github.com/michaeldcanady/go-onedrive/internal/feature/fs"
	"github.com/michaeldcanady/go-onedrive/internal/feature/logger"
)

// Handler executes the drive download operation.
type Handler struct {
	manager shared.Service
	log     logger.Logger
}

// NewHandler initializes a new instance of the drive download Handler.
func NewHandler(m shared.Service, l logger.Logger) *Handler {
	return &Handler{
		manager: m,
		log:     l,
	}
}

// Handle downloads an item from the remote filesystem to the local destination.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx).With(
		logger.String("source", opts.Source),
		logger.String("destination", opts.Destination),
	)

	log.Info("starting download")

	// Ensure destination has local: prefix if no prefix is present.
	// TODO: what if user provides onedrive: prefix for destination?
	destination := opts.Destination
	if !strings.Contains(destination, ":") {
		log.Debug("adding local prefix to destination")
		destination = "local:" + destination
	}

	cpOpts := shared.CopyOptions{
		Recursive: opts.Recursive,
		Overwrite: true,
	}

	log.Debug("delegating to filesystem manager for copy")
	if err := h.manager.Copy(ctx, opts.Source, destination, cpOpts); err != nil {
		log.Error("download failed", logger.Error(err))
		return fmt.Errorf("failed to download %s to %s: %w", opts.Source, destination, err)
	}

	log.Info("download completed successfully")
	fmt.Fprintf(opts.Stdout, "Downloaded %s to %s\n", opts.Source, opts.Destination)
	return nil
}
