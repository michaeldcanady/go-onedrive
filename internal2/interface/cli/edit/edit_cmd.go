package edit

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

// EditCmd handles the execution logic for the 'edit' command.
// It coordinates downloading a file, launching an editor, and uploading changes.
type EditCmd struct {
	container di.Container
	logger    infralogging.Logger
}

// NewEditCmd creates a new EditCmd instance with the provided dependency container.
func NewEditCmd(container di.Container) *EditCmd {
	return &EditCmd{
		container: container,
	}
}

// WithLogger allows injecting a logger into EditCmd.
func (c *EditCmd) WithLogger(logger infralogging.Logger) *EditCmd {
	c.logger = logger
	return c
}

// Run executes the edit lifecycle.
// 1. Ensures a logger is available.
// 2. Fetches the file content from OneDrive.
// 3. Launches the local editor with a temporary file.
// 4. Checks for changes using SHA-256 hashes.
// 5. Uploads the updated file if changes were made.
func (c *EditCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if ctx == nil {
		ctx = context.Background()
	}

	if c.logger == nil {
		logger, err := util.EnsureLogger(c.container, loggerID)
		if err != nil {
			return util.NewCommandErrorWithNameWithError(commandName, err)
		}
		c.logger = logger
	}

	c.logger.Info("starting edit command", infralogging.String("path", opts.Path))

	fsSvc := c.container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(commandName, "filesystem service is nil")
	}

	// 1. Read File
	reader, err := fsSvc.ReadFile(ctx, opts.Path, domainfs.ReadOptions{})
	if err != nil {
		return util.NewCommandError(commandName, "failed to read file from OneDrive", err)
	}
	defer reader.Close()

	// 2. Track Original Hash
	origHash := sha256.New()
	name := Name(opts.Path)
	ext := filepath.Ext(opts.Path)

	// 3. Launch Editor
	editorSvc := NewEditorService(c.container.EnvironmentService(), c.logger).
		WithIO(opts.Stdin, opts.Stdout, opts.Stderr)

	editedBytes, tmpPath, err := editorSvc.LaunchTempFile(fmt.Sprintf("%s-edit-", name), ext, io.TeeReader(reader, origHash))
	defer os.Remove(tmpPath)
	if err != nil {
		return util.NewCommandErrorWithNameWithError(commandName, err)
	}

	// 4. Compare Hashes
	origHashSum := hex.EncodeToString(origHash.Sum(nil))
	editedHash := sha256.Sum256(editedBytes)

	if origHashSum == hex.EncodeToString(editedHash[:]) {
		c.logger.Info("no changes detected, skipping upload")
		fmt.Fprintln(opts.Stdout, "No changes detected.")
		return nil
	}

	// 5. Upload Changes
	c.logger.Info("changes detected, uploading updated file")
	_, err = fsSvc.WriteFile(ctx, opts.Path, bytes.NewReader(editedBytes), domainfs.WriteOptions{Overwrite: opts.Force})
	if err != nil {
		if err == domainfs.ErrPrecondition {
			return util.NewCommandErrorWithNameWithMessage(commandName, "failed to upload: the file has been modified in the cloud. Use --force to overwrite anyway.")
		}
		return util.NewCommandError(commandName, "failed to upload updated file", err)
	}

	c.logger.Info("file updated successfully",
		infralogging.Duration("duration", time.Since(start)),
	)
	fmt.Fprintf(opts.Stdout, "File %q updated successfully.\n", opts.Path)

	return nil
}

