package edit

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/michaeldcanady/go-onedrive/internal/errors"
	registry "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/fs/editor"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/pkg/fsm"
)

// Handler executes the drive edit operation.
type Handler struct {
	fs     registry.Service
	editor editor.Service
	log    logger.Logger
}

// editContext holds the shared data for the edit state machine.
type editContext struct {
	handler *Handler
	opts    Options
	item    registry.Item
	content io.ReadCloser
	newData []byte
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
	data := &editContext{
		handler: h,
		opts:    opts,
	}

	machine := fsm.NewMachine(data)
	return machine.Run(ctx, fsm.StateFunc[editContext](fetchMetadataState))
}

// fetchMetadataState retrieves the file metadata to prepare for editing.
func fetchMetadataState(ctx context.Context, data *editContext) (fsm.State[editContext], error) {
	log := data.handler.log.WithContext(ctx).With(logger.String("path", data.opts.Path.String()))
	log.Debug("fetching metadata")

	item, err := data.handler.fs.Get(ctx, data.opts.Path)
	if err != nil {
		wrapped := cli.WrapError(err, data.opts.Path.String())
		data.handler.log.Error(wrapped.Error(), errors.LogFields(wrapped)...)
		return nil, wrapped
	}

	data.item = item
	return fsm.StateFunc[editContext](downloadState), nil
}

// downloadState reads the file content for editing.
func downloadState(ctx context.Context, data *editContext) (fsm.State[editContext], error) {
	log := data.handler.log.WithContext(ctx).With(logger.String("path", data.opts.Path.String()))
	log.Debug("reading file for editing")

	r, err := data.handler.fs.ReadFile(ctx, data.opts.Path, registry.ReadOptions{})
	if err != nil {
		wrapped := cli.WrapError(err, data.opts.Path.String())
		data.handler.log.Error(wrapped.Error(), errors.LogFields(wrapped)...)
		return nil, wrapped
	}

	data.content = r
	return fsm.StateFunc[editContext](editState), nil
}

// editState launches the external editor with the file content and captures the edited data.
func editState(ctx context.Context, data *editContext) (fsm.State[editContext], error) {
	log := data.handler.log.WithContext(ctx).With(logger.String("path", data.opts.Path.String()))
	log.Debug("launching external editor")

	defer data.content.Close()

	newData, _, err := data.handler.editor.LaunchTempFile(ctx, "odc-edit", ".txt", data.content)
	if err != nil {
		wrapped := cli.WrapError(err, data.opts.Path.String())
		data.handler.log.Error(wrapped.Error(), errors.LogFields(wrapped)...)
		return nil, wrapped
	}

	if newData == nil {
		log.Info("no changes detected, skipping upload")
		_, _ = fmt.Fprintln(data.opts.Stdout, "No changes detected.")
		return nil, nil
	}

	data.newData = newData
	return fsm.StateFunc[editContext](uploadState), nil
}

// uploadState uploads the edited content back to the file system, handling conflicts if they arise.
func uploadState(ctx context.Context, data *editContext) (fsm.State[editContext], error) {
	log := data.handler.log.WithContext(ctx).With(logger.String("path", data.opts.Path.String()))
	log.Info("uploading changes", logger.Int("size", len(data.newData)))

	writeOpts := registry.WriteOptions{
		Overwrite: data.opts.Force,
	}
	// Use ETag for optimistic concurrency control if not forcing
	if !data.opts.Force {
		writeOpts.IfMatch = data.item.ETag
	}

	_, err := data.handler.fs.WriteFile(ctx, data.opts.Path, bytes.NewReader(data.newData), writeOpts)
	if err != nil {
		if errors.Is(err, errors.CodePrecondition) {
			log.Warn("conflict detected, upload aborted")
			return nil, errors.NewAppError(errors.CodeConflict, err, "conflict detected during upload", "Use --force to overwrite changes.")
		}
		wrapped := cli.WrapError(err, data.opts.Path.String())
		data.handler.log.Error(wrapped.Error(), errors.LogFields(wrapped)...)
		return nil, wrapped
	}

	log.Info("edit completed successfully")
	_, _ = fmt.Fprintf(data.opts.Stdout, "successfully updated file \"%s\"\n", data.opts.Path.String())
	return nil, nil
}
