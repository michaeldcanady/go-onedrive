package cat

import (
	"context"
	"fmt"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/registry"
	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/shared"
	"github.com/michaeldcanady/go-onedrive/internal2/core/logger"
)

// Handler executes the drive cat operation.
type Handler struct {
	fs  registry.Service
	log logger.Logger
}

// NewHandler initializes a new instance of the drive cat Handler.
func NewHandler(fs registry.Service, l logger.Logger) *Handler {
	return &Handler{
		fs:  fs,
		log: l,
	}
}

// Handle retrieves and writes the content of a file to stdout.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	h.log.Info("reading file content", logger.String("path", opts.Path))

	provider, path, err := h.fs.Resolve(ctx, opts.Path)
	if err != nil {
		return fmt.Errorf("failed to resolve path %s: %w", opts.Path, err)
	}

	reader, err := provider.ReadFile(ctx, path, shared.ReadOptions{})
	if err != nil {
		return fmt.Errorf("failed to read file at %s: %w", path, err)
	}
	defer reader.Close()

	if _, err := io.Copy(opts.Stdout, reader); err != nil {
		return fmt.Errorf("failed to write content to output: %w", err)
	}

	return nil
}
