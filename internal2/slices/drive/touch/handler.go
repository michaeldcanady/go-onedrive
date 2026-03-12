package touch

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/registry"
	"github.com/michaeldcanady/go-onedrive/internal2/core/logger"
)

// Handler executes the drive touch operation.
type Handler struct {
	fs  registry.Service
	log logger.Logger
}

// NewHandler initializes a new instance of the drive touch Handler.
func NewHandler(fs registry.Service, l logger.Logger) *Handler {
	return &Handler{
		fs:  fs,
		log: l,
	}
}

// Handle creates an empty file or updates the timestamp of an existing one.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("touching file", logger.String("path", opts.Path))

	provider, path, err := h.fs.Resolve(ctx, opts.Path)
	if err != nil {
		return fmt.Errorf("failed to resolve path %s: %w", opts.Path, err)
	}

	if _, err := provider.Touch(ctx, path); err != nil {
		return fmt.Errorf("failed to touch file at %s: %w", path, err)
	}

	fmt.Fprintf(opts.Stdout, "File touched: %s\n", opts.Path)
	return nil
}
