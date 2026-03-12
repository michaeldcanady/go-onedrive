package mkdir

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs/registry"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
)

// Handler executes the drive mkdir operation.
type Handler struct {
	fs  registry.Service
	log logger.Logger
}

// NewHandler initializes a new instance of the drive mkdir Handler.
func NewHandler(fs registry.Service, l logger.Logger) *Handler {
	return &Handler{
		fs:  fs,
		log: l,
	}
}

// Handle creates a new directory.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("creating directory", logger.String("path", opts.Path))

	provider, path, err := h.fs.Resolve(ctx, opts.Path)
	if err != nil {
		return fmt.Errorf("failed to resolve path %s: %w", opts.Path, err)
	}

	if err := provider.Mkdir(ctx, path); err != nil {
		return fmt.Errorf("failed to create directory at %s: %w", path, err)
	}

	fmt.Fprintf(opts.Stdout, "Directory created: %s\n", opts.Path)
	return nil
}
