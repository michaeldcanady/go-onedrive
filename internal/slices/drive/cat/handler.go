package cat

import (
	"context"
	"fmt"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/core/fs/registry"
	"github.com/michaeldcanady/go-onedrive/internal/core/fs/shared"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
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
	log := h.log.WithContext(ctx).With(logger.String("path", opts.Path))

	log.Info("reading file content")

	log.Debug("resolving path provider")
	provider, subPath, err := h.fs.Resolve(ctx, opts.Path)
	if err != nil {
		log.Error("failed to resolve path", logger.Error(err))
		return fmt.Errorf("failed to resolve path %s: %w", opts.Path, err)
	}
	log.Debug("resolved path provider", logger.String("provider", provider.Name()))

	log.Debug("fetching file reader")
	reader, err := provider.ReadFile(ctx, subPath, shared.ReadOptions{})
	if err != nil {
		log.Error("failed to read file", logger.Error(err))
		return fmt.Errorf("failed to read file at %s: %w", opts.Path, err)
	}
	defer reader.Close()

	log.Debug("copying content to stdout")
	if _, err := io.Copy(opts.Stdout, reader); err != nil {
		log.Error("failed to write output", logger.Error(err))
		return fmt.Errorf("failed to write content to output: %w", err)
	}

	log.Info("cat completed successfully")
	return nil
}
