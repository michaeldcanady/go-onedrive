package upload

import (
	"context"
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs/shared"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
)

// Handler executes the drive upload operation.
type Handler struct {
	manager shared.Service
	log     logger.Logger
}

// NewHandler initializes a new instance of the drive upload Handler.
func NewHandler(m shared.Service, l logger.Logger) *Handler {
	return &Handler{
		manager: m,
		log:     l,
	}
}

// Handle uploads an item from the local filesystem to the remote destination.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("uploading item",
		logger.String("source", opts.Source),
		logger.String("destination", opts.Destination),
		logger.Bool("recursive", opts.Recursive),
	)

	// Ensure source has local: prefix if no prefix is present.
	source := opts.Source
	if !strings.Contains(source, ":") {
		source = "local:" + source
	}

	cpOpts := shared.CopyOptions{
		Recursive: opts.Recursive,
		Overwrite: true,
	}

	if err := h.manager.Copy(ctx, source, opts.Destination, cpOpts); err != nil {
		return fmt.Errorf("failed to upload %s to %s: %w", opts.Source, opts.Destination, err)
	}

	fmt.Fprintf(opts.Stdout, "Uploaded %s to %s\n", opts.Source, opts.Destination)
	return nil
}
