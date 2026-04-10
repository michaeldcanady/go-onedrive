package edit

import (
	"bytes"
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
	registry "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/fs/editor"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Handler executes the drive edit operation.
type Handler struct {
	fs     registry.Service
	editor editor.Service
	log    logger.Logger
}

// NewHandler initializes a new instance of the drive edit Handler.
func NewHandler(
	fs registry.Service,
	editor editor.Service,
	l logger.Service,
) *Handler {
	cliLog, _ := l.CreateLogger("drive-edit")
	return &Handler{
		fs:     fs,
		editor: editor,
		log:    cliLog,
	}
}

// Handle retrieves a file, opens it in an editor, and uploads changes.
func (h *Handler) Handle(ctx context.Context, opts Options) error {
	log := h.log.WithContext(ctx).With(logger.String("path", opts.Path))

	log.Info("editing file")

	log.Debug("fetching metadata")
	item, err := h.fs.Get(ctx, opts.Path)
	if err != nil {
		wrapped := cli.WrapError(err, opts.Path)
		h.log.Error(wrapped.Error(), errors.LogFields(wrapped)...)
		return wrapped
	}

	log.Debug("reading file for editing")
	r, err := h.fs.ReadFile(ctx, opts.Path, registry.ReadOptions{})
	if err != nil {
		wrapped := cli.WrapError(err, opts.Path)
		h.log.Error(wrapped.Error(), errors.LogFields(wrapped)...)
		return wrapped
	}
	defer r.Close()

	log.Debug("launching external editor")
	newData, _, err := h.editor.LaunchTempFile(ctx, "odc-edit", ".txt", r)
	if err != nil {
		wrapped := cli.WrapError(err, opts.Path)
		h.log.Error(wrapped.Error(), errors.LogFields(wrapped)...)
		return wrapped
	}

	if newData == nil {
		log.Info("no changes detected, skipping upload")
		_, _ = fmt.Fprintln(opts.Stdout, "No changes detected.")
		return nil
	}

	log.Info("uploading changes", logger.Int("size", len(newData)))
	writeOpts := registry.WriteOptions{
		Overwrite: opts.Force,
	}
	// Use ETag for optimistic concurrency control if not forcing
	if !opts.Force {
		writeOpts.IfMatch = item.ETag
	}

	_, err = h.fs.WriteFile(ctx, opts.Path, bytes.NewReader(newData), writeOpts)
	if err != nil {
		if errors.Is(err, errors.CodePrecondition) {
			log.Warn("conflict detected, upload aborted")
			return errors.NewAppError(errors.CodeConflict, err, "conflict detected during upload", "Use --force to overwrite changes.")
		}
		wrapped := cli.WrapError(err, opts.Path)
		h.log.Error(wrapped.Error(), errors.LogFields(wrapped)...)
		return wrapped
	}

	log.Info("edit completed successfully")
	_, _ = fmt.Fprintf(opts.Stdout, "successfully updated file \"%s\"\n", opts.Path)
	return nil
}
