package cat

import (
	"context"
	"fmt"
	"io"

	registry "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
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
	log := h.log.WithContext(ctx).With(logger.String("path", opts.URI.String()))

	log.Info("reading file content")

	log.Debug("fetching file reader")
	reader, err := h.fs.ReadFile(ctx, opts.URI.String(), registry.ReadOptions{})
	if err != nil {
		log.Error("failed to read file", logger.Error(err))
		return fmt.Errorf("failed to read file at %s: %w", opts.URI.String(), err)
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
