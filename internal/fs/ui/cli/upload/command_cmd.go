package upload

import (
	"context"
	"fmt"
	"strings"

	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
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
	log := h.log.WithContext(ctx).With(
		logger.String("source", opts.Source),
		logger.String("destination", opts.Destination),
	)

	log.Info("starting upload")

	// Ensure source has local: prefix if no prefix is present.
	source := opts.Source
	if !strings.Contains(source, ":") {
		log.Debug("adding local prefix to source")
		source = "local:" + source
	}

	cpOpts := shared.CopyOptions{
		Recursive: opts.Recursive,
		Overwrite: true,
	}

	log.Debug("delegating to filesystem manager for copy")
	if err := h.manager.Copy(ctx, source, opts.Destination, cpOpts); err != nil {
		log.Error("upload failed", logger.Error(err))
		return fmt.Errorf("failed to upload %s to %s: %w", opts.Source, opts.Destination, err)
	}

	log.Info("upload completed successfully")
	fmt.Fprintf(opts.Stdout, "Uploaded %s to %s\n", opts.Source, opts.Destination)
	return nil
}
