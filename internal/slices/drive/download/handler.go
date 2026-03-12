package download

import (
	"context"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs/shared"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
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
	h.log.Info("downloading item",
		logger.String("source", opts.Source),
		logger.String("destination", opts.Destination),
		logger.Bool("recursive", opts.Recursive),
	)

	// Ensure destination has local: prefix if no prefix is present.
	destination := opts.Destination
	if !strings.Contains(destination, ":") {
		destination = "local:" + destination
	}

	cpOpts := shared.CopyOptions{
		Recursive: opts.Recursive,
		Overwrite: true,
	}

	if err := h.manager.Copy(ctx, opts.Source, destination, cpOpts); err != nil {
		return fmt.Errorf("failed to download %s to %s: %w", opts.Source, destination, err)
	}

	fmt.Fprintf(opts.Stdout, "Downloaded %s to %s\n", opts.Source, opts.Destination)
	return nil
}
