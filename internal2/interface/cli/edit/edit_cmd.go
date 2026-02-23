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

	c.logger.Debug("resolving filesystem service")
	fsSvc := c.container.FS()
	if fsSvc == nil {
		c.logger.Error("filesystem service is nil")
		return util.NewCommandErrorWithNameWithMessage(commandName, "filesystem service is nil")
	}

	// 1. Read File
	c.logger.Debug("reading file from OneDrive", infralogging.String("path", opts.Path))
	reader, err := fsSvc.ReadFile(ctx, opts.Path, domainfs.ReadOptions{})
	if err != nil {
		c.logger.Error("failed to read file from OneDrive", 
			infralogging.String("path", opts.Path),
			infralogging.Error(err))
		return util.NewCommandError(commandName, "failed to read file from OneDrive", err)
	}
	defer reader.Close()

	// 2. Track Original Hash
	origHash := sha256.New()
	name := Name(opts.Path)
	ext := filepath.Ext(opts.Path)

	// 3. Launch Editor
	c.logger.Debug("resolving editor service")
	editorSvc := NewEditorService(c.container.EnvironmentService(), c.logger).
		WithIO(opts.Stdin, opts.Stdout, opts.Stderr)

	c.logger.Info("launching local editor", 
		infralogging.String("file", opts.Path),
		infralogging.String("extension", ext))
	editedBytes, tmpPath, err := editorSvc.LaunchTempFile(fmt.Sprintf("%s-edit-", name), ext, io.TeeReader(reader, origHash))
	if tmpPath != "" {
		defer func() {
			c.logger.Debug("removing temporary file", infralogging.String("path", tmpPath))
			os.Remove(tmpPath)
		}()
	}
	if err != nil {
		c.logger.Error("editor launch or execution failed", infralogging.Error(err))
		return util.NewCommandErrorWithNameWithError(commandName, err)
	}

	// 4. Compare Hashes
	c.logger.Debug("calculating hashes for change detection")
	origHashSum := hex.EncodeToString(origHash.Sum(nil))
	editedHash := sha256.Sum256(editedBytes)
	editedHashSum := hex.EncodeToString(editedHash[:])

	c.logger.Debug("hash comparison", 
		infralogging.String("original", origHashSum),
		infralogging.String("edited", editedHashSum))

	if origHashSum == editedHashSum {
		c.logger.Info("no changes detected, skipping upload")
		fmt.Fprintln(opts.Stdout, "No changes detected.")
		return nil
	}

	// 5. Upload Changes
	c.logger.Info("changes detected, uploading updated file", 
		infralogging.String("path", opts.Path),
		infralogging.Bool("force", opts.Force))
	_, err = fsSvc.WriteFile(ctx, opts.Path, bytes.NewReader(editedBytes), domainfs.WriteOptions{Overwrite: opts.Force})
	if err != nil {
		if err == domainfs.ErrPrecondition {
			c.logger.Warn("upload rejected: file modified in cloud", infralogging.String("path", opts.Path))
			return util.NewCommandErrorWithNameWithMessage(commandName, "failed to upload: the file has been modified in the cloud. Use --force to overwrite anyway.")
		}
		c.logger.Error("failed to upload updated file", 
			infralogging.String("path", opts.Path),
			infralogging.Error(err))
		return util.NewCommandError(commandName, "failed to upload updated file", err)
	}

	c.logger.Info("file updated successfully",
		infralogging.String("path", opts.Path),
		infralogging.Duration("duration", time.Since(start)),
	)
	fmt.Fprintf(opts.Stdout, "File %q updated successfully.\n", opts.Path)

	return nil
}

