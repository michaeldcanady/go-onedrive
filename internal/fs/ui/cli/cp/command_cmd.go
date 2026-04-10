package cp

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Handler executes the drive cp operation.
type Handler struct {
	manager shared.Service
	log     logger.Logger
}

// NewHandler initializes a new instance of the drive cp Handler.
func NewHandler(
	m shared.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("drive-cp")
	return &Handler{
		manager: m,
		log:     cliLog,
	}
}

// Handle copies an item from the source to the destination.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx).With(
		logger.String("source", opts.Source),
		logger.String("destination", opts.Destination),
		logger.Bool("recursive", opts.Recursive),
	)

	log.Info("starting copy operation")

	cpOpts := shared.CopyOptions{
		Recursive: opts.Recursive,
		Overwrite: true, // Default to overwrite for now, can be a flag later.
	}

	log.Debug("delegating to filesystem manager for copy")
	if err := h.manager.Copy(ctx, opts.Source, opts.Destination, cpOpts); err != nil {
		wrapped := cli.WrapError(err, opts.Source)
		h.log.Error(wrapped.Error(), errors.LogFields(wrapped)...)
		return wrapped
	}

	log.Info("copy completed successfully")
	fmt.Fprintf(opts.Stdout, "Copied %s to %s\n", opts.Source, opts.Destination)
	return nil
}
