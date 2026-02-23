package edit

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	domaineditor "github.com/michaeldcanady/go-onedrive/internal2/domain/editor"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

// EditCmd handles the execution logic for the 'edit' command.
// It coordinates downloading a file, launching an editor, and uploading changes.
type EditCmd struct {
	container di.Container
	logger    infralogging.Logger
	editor    domaineditor.Service
}

// NewEditCmd creates a new EditCmd instance with the provided dependency container.
func NewEditCmd(container di.Container) *EditCmd {
	return &EditCmd{
		container: container,
	}
}

// WithEditor allows injecting a custom editor into EditCmd.
func (c *EditCmd) WithEditor(editor domaineditor.Service) *EditCmd {
	c.editor = editor
	return c
}

// WithLogger allows injecting a logger into EditCmd.
func (c *EditCmd) WithLogger(logger infralogging.Logger) *EditCmd {
	c.logger = logger
	return c
}

// Run executes the edit lifecycle.
func (c *EditCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.ensureDependencies(ctx); err != nil {
		return err
	}

	c.logger.Info("starting edit command", infralogging.String("path", opts.Path))

	// 1. Download
	reader, err := c.downloadFile(ctx, opts.Path)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 2. Edit
	origHash := sha256.New()
	name := Name(opts.Path)
	ext := filepath.Ext(opts.Path)

	editedBytes, tmpPath, err := c.launchEditor(opts, name, ext, io.TeeReader(reader, origHash))
	if tmpPath != "" {
		defer os.Remove(tmpPath)
	}
	if err != nil {
		return err
	}

	// 3. Compare and Upload
	if c.hasChanges(origHash, editedBytes) {
		if err := c.uploadChanges(ctx, opts.Path, editedBytes, opts.Force); err != nil {
			return err
		}
		c.logger.Info("file updated successfully",
			infralogging.String("path", opts.Path),
			infralogging.Duration("duration", time.Since(start)),
		)
		fmt.Fprintf(opts.Stdout, "File %q updated successfully.\n", opts.Path)
	} else {
		c.logger.Info("no changes detected, skipping upload")
		fmt.Fprintln(opts.Stdout, "No changes detected.")
	}

	return nil
}

func (c *EditCmd) ensureDependencies(ctx context.Context) error {
	if c.logger == nil {
		logger, err := util.EnsureLogger(c.container, loggerID)
		if err != nil {
			return util.NewCommandErrorWithNameWithError(commandName, err)
		}
		c.logger = logger
	}

	if c.editor == nil {
		c.editor = c.container.Editor()
	}

	c.logger = c.logger.WithContext(ctx).With(infralogging.String("correlationID", util.CorrelationIDFromContext(ctx)))
	return nil
}

func (c *EditCmd) downloadFile(ctx context.Context, path string) (io.ReadCloser, error) {
	fsSvc := c.container.FS()
	if fsSvc == nil {
		return nil, util.NewCommandErrorWithNameWithMessage(commandName, "filesystem service is nil")
	}

	c.logger.Debug("reading file from OneDrive", infralogging.String("path", path))
	reader, err := fsSvc.ReadFile(ctx, path, domainfs.ReadOptions{})
	if err != nil {
		c.logger.Error("failed to read file from OneDrive",
			infralogging.String("path", path),
			infralogging.Error(err))
		return nil, util.NewCommandError(commandName, "failed to read file from OneDrive", err)
	}
	return reader, nil
}

func (c *EditCmd) launchEditor(opts Options, name, ext string, reader io.Reader) ([]byte, string, error) {
	c.logger.Debug("launching local editor",
		infralogging.String("file", opts.Path),
		infralogging.String("extension", ext))

	c.editor.WithIO(opts.Stdin, opts.Stdout, opts.Stderr)

	editedBytes, tmpPath, err := c.editor.LaunchTempFile(fmt.Sprintf("%s-edit-", name), ext, reader)
	if err != nil {
		c.logger.Error("editor launch or execution failed", infralogging.Error(err))
		return nil, tmpPath, util.NewCommandErrorWithNameWithError(commandName, err)
	}
	return editedBytes, tmpPath, nil
}

func (c *EditCmd) hasChanges(origHash hash.Hash, editedBytes []byte) bool {
	origHashSum := hex.EncodeToString(origHash.Sum(nil))
	editedHash := sha256.Sum256(editedBytes)
	editedHashSum := hex.EncodeToString(editedHash[:])

	c.logger.Debug("hash comparison",
		infralogging.String("original", origHashSum),
		infralogging.String("edited", editedHashSum))

	return origHashSum != editedHashSum
}

func (c *EditCmd) uploadChanges(ctx context.Context, path string, content []byte, force bool) error {
	fsSvc := c.container.FS()
	c.logger.Info("uploading updated file",
		infralogging.String("path", path),
		infralogging.Bool("force", force))

	_, err := fsSvc.WriteFile(ctx, path, bytes.NewReader(content), domainfs.WriteOptions{Overwrite: force})
	if err != nil {
		c.logger.Error("failed to upload updated file",
			infralogging.String("path", path),
			infralogging.Error(err))
		return util.NewCommandError(commandName, "failed to upload updated file", err)
	}
	return nil
}
