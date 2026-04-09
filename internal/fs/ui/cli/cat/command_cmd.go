package cat

import (
	"context"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
	registry "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Handler executes the drive cat operation.
type Handler struct {
	fs  registry.Service
	log logger.Logger
}

// NewHandler initializes a new instance of the drive cat Handler.
func NewHandler(
	fs registry.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("drive-cat")
	return &Handler{
		fs:  fs,
		log: cliLog,
	}
}

// Handle retrieves and writes the content of a file to stdout.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx).With(logger.String("path", opts.Path))

	log.Info("reading file content")

	log.Debug("fetching file reader")
	reader, err := h.fs.ReadFile(ctx, opts.Path, registry.ReadOptions{})
	if err != nil {
		h.log.Error(err.Error(), errors.LogFields(err)...)
		return err
	}
	defer reader.Close()

	log.Debug("copying content to stdout")
	if _, err := io.Copy(opts.Stdout, reader); err != nil {
		h.log.Error(err.Error(), errors.LogFields(err)...)
		return err
	}

	log.Info("cat completed successfully")
	return nil
}
