package download

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Handler executes the drive download operation.
type Handler struct {
	manager shared.Service
	log     logger.Logger
}

// NewHandler initializes a new instance of the drive download Handler.
func NewHandler(
	m shared.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("drive-download")
	return &Handler{
		manager: m,
		log:     cliLog,
	}
}

// Handle downloads an item from the remote filesystem to the local destination.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx).With(
		logger.String("source", opts.Source.String()),
		logger.String("destination", opts.Destination.String()),
	)

	log.Info("starting download")

	// Ensure destination has local: prefix if no prefix is present.
	// TODO: what if user provides onedrive: prefix for destination?
	destination := opts.Destination
	if destination.Provider == "" {
		log.Debug("setting local provider for destination")
		destination.Provider = "local"
	}

	cpOpts := shared.CopyOptions{
		Recursive: opts.Recursive,
		Overwrite: true,
	}

	log.Debug("delegating to filesystem manager for copy")
	if err := h.manager.Copy(ctx, opts.Source, destination, cpOpts); err != nil {
		wrapped := cli.WrapError(err, opts.Source.String())
		h.log.Error(wrapped.Error(), errors.LogFields(wrapped)...)
		return wrapped
	}

	log.Info("download completed successfully")
	fmt.Fprintf(opts.Stdout, "Downloaded %s to %s\n", opts.Source.String(), opts.Destination.String())
	return nil
}
