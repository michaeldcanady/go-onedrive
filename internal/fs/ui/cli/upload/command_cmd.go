package upload

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/errors"
	shared "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/pkg/events"
)

// Handler executes the drive upload operation.
type Handler struct {
	container di.Container
	manager   shared.Service
	log       logger.Logger
}

// NewHandler initializes a new instance of the drive upload Handler.
func NewHandler(
	c di.Container,
) *Handler {
	cliLog, _ := c.Logger().CreateLogger("drive-upload")
	return &Handler{
		container: c,
		manager:   c.FS(),
		log:       cliLog,
	}
}

// Handle uploads an item from the local filesystem to the remote destination.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx).With(
		logger.String("source", opts.Source.String()),
		logger.String("destination", opts.Destination.String()),
	)

	log.Info("starting upload")

	// Subscribe to progress events
	h.container.Events().Subscribe((shared.FileTransferProgressEvent{}).Name(), func(e events.Event) {
		if pe, ok := e.(shared.FileTransferProgressEvent); ok {
			percent := float64(pe.Transferred) / float64(pe.Total) * 100
			fmt.Fprintf(opts.Stdout, "\rUploading %s: %.2f%% (%d/%d bytes)", pe.Path, percent, pe.Transferred, pe.Total)
		}
	})

	// Ensure source has local: prefix if no prefix is present.
	source := opts.Source
	if source.Provider == "" {
		log.Debug("setting local provider for source")
		source.Provider = "local"
	}

	cpOpts := shared.CopyOptions{
		Recursive: opts.Recursive,
		Overwrite: true,
	}

	log.Debug("delegating to filesystem manager for copy")
	if err := h.manager.Copy(ctx, source, opts.Destination, cpOpts); err != nil {
		wrapped := cli.WrapError(err, source.String())
		h.log.Error(wrapped.Error(), errors.LogFields(wrapped)...)
		return wrapped
	}

	log.Info("upload completed successfully")
	fmt.Fprintf(opts.Stdout, "\nUploaded %s to %s\n", opts.Source, opts.Destination)
	return nil
}
