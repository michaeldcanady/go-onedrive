package edit

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	domaineditor "github.com/michaeldcanady/go-onedrive/internal2/domain/editor"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	infrafile "github.com/michaeldcanady/go-onedrive/internal2/infra/file"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

// ConflictHandler defines the interface for resolving file upload conflicts.
type ConflictHandler interface {
	HandleConflict(ctx context.Context, path string, content []byte, tmpPath string) (bool, string, error)
}

// DefaultConflictHandler implements ConflictHandler using promptui for user interaction.
type DefaultConflictHandler struct {
	cmd    *EditCmd
	stdout io.Writer
	stderr io.Writer
}

func (h *DefaultConflictHandler) HandleConflict(ctx context.Context, path string, content []byte, tmpPath string) (bool, string, error) {
	prompt := promptui.Select{
		Label: "Conflict: The file has been modified in the cloud. How would you like to proceed?",
		Items: []string{
			"Overwrite (Force Upload)",
			"Save as Copy",
			"Keep Local & Abort",
			"Discard & Abort",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return false, "", util.NewCommandErrorWithNameWithError(commandName, err)
	}

	switch result {
	case "Overwrite (Force Upload)":
		if err := h.cmd.uploadChanges(ctx, path, content, true); err != nil {
			fmt.Fprintf(h.stderr, "\nUpload failed. Your changes are saved locally at: %s\n", tmpPath)
			return false, "", err
		}
		return true, path, nil
	case "Save as Copy":
		ext := filepath.Ext(path)
		base := strings.TrimSuffix(path, ext)
		copyPath := fmt.Sprintf("%s-copy%s", base, ext)

		// Ask for confirmation/edit of the copy path
		namePrompt := promptui.Prompt{
			Label:   "Save copy as",
			Default: copyPath,
		}
		finalPath, err := namePrompt.Run()
		if err != nil {
			return false, "", util.NewCommandErrorWithNameWithError(commandName, err)
		}

		if err := h.cmd.uploadChanges(ctx, finalPath, content, false); err != nil {
			fmt.Fprintf(h.stderr, "\nUpload failed. Your changes are saved locally at: %s\n", tmpPath)
			return false, "", err
		}
		return true, finalPath, nil

	case "Keep Local & Abort":
		fmt.Fprintf(h.stderr, "Aborted. Your changes are saved locally at: %s\n", tmpPath)
		return false, "", nil

	case "Discard & Abort":
		return true, "", nil

	default:
		return false, "", util.NewCommandErrorWithNameWithMessage(commandName, "invalid selection")
	}
}

// EditCmd handles the execution logic for the 'edit' command.
// It coordinates downloading a file, launching an editor, and uploading changes.
type EditCmd struct {
	container       di.Container
	logger          infralogging.Logger
	editor          domaineditor.Service
	conflictHandler ConflictHandler
	etag            string
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

// WithConflictHandler allows injecting a custom conflict handler into EditCmd.
func (c *EditCmd) WithConflictHandler(handler ConflictHandler) *EditCmd {
	c.conflictHandler = handler
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
	reader, etag, err := c.downloadFile(ctx, opts.Path)
	if err != nil {
		return err
	}
	defer reader.Close()
	c.etag = etag

	// 2. Edit
	origHash := sha256.New()
	name := Name(opts.Path)
	ext := filepath.Ext(opts.Path)

	editedBytes, tmpPath, err := c.launchEditor(opts, name, ext, io.TeeReader(reader, origHash))
	// Always remove the temp file if there's no error, but keep it if upload fails later
	shouldRemoveTemp := true
	defer func() {
		if shouldRemoveTemp && tmpPath != "" {
			c.logger.Debug("removing temporary file", infralogging.String("path", tmpPath))
			os.Remove(tmpPath)
		}
	}()

	if err != nil {
		return err
	}

	// 3. Compare and Upload
	if c.hasChanges(origHash, editedBytes) {
		var finalPath string
		err := c.uploadChanges(ctx, opts.Path, editedBytes, opts.Force)
		if err == nil {
			finalPath = opts.Path
		} else {
			if errors.Is(err, infrafile.ErrPrecondition) && !opts.Force {
				c.logger.Warn("conflict detected, initiating interactive resolution")
				shouldRemoveTemp, finalPath, err = c.conflictHandler.HandleConflict(ctx, opts.Path, editedBytes, tmpPath)
				if err != nil {
					return err
				}
			} else {
				shouldRemoveTemp = false
				fmt.Fprintf(opts.Stderr, "\nUpload failed. Your changes are saved locally at: %s\n", tmpPath)
				return err
			}
		}

		if finalPath != "" {
			c.logger.Info("file updated successfully",
				infralogging.String("path", finalPath),
				infralogging.Duration("duration", time.Since(start)),
			)
			fmt.Fprintf(opts.Stdout, "File %q updated successfully.\n", finalPath)
		}
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

	if c.conflictHandler == nil {
		c.conflictHandler = &DefaultConflictHandler{cmd: c, stdout: os.Stdout, stderr: os.Stderr}
	}

	c.logger = c.logger.WithContext(ctx).With(infralogging.String("correlationID", util.CorrelationIDFromContext(ctx)))
	return nil
}

func (c *EditCmd) downloadFile(ctx context.Context, path string) (io.ReadCloser, string, error) {
	fsSvc := c.container.FS()
	if fsSvc == nil {
		return nil, "", util.NewCommandErrorWithNameWithMessage(commandName, "filesystem service is nil")
	}

	c.logger.Debug("retrieving file metadata", infralogging.String("path", path))
	item, err := fsSvc.Get(ctx, path)
	if err != nil {
		c.logger.Error("failed to get file metadata",
			infralogging.String("path", path),
			infralogging.Error(err))
		return nil, "", util.NewCommandError(commandName, "failed to get file metadata", err)
	}

	c.logger.Debug("reading file from OneDrive", infralogging.String("path", path))
	reader, err := fsSvc.ReadFile(ctx, path, domainfs.ReadOptions{})
	if err != nil {
		c.logger.Error("failed to read file from OneDrive",
			infralogging.String("path", path),
			infralogging.Error(err))
		return nil, "", util.NewCommandError(commandName, "failed to read file from OneDrive", err)
	}
	return reader, item.ETag, nil
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

	ifMatch := c.etag
	if force {
		ifMatch = ""
	}

	c.logger.Info("uploading updated file",
		infralogging.String("path", path),
		infralogging.Bool("force", force),
		infralogging.String("ifMatch", ifMatch))

	_, err := fsSvc.WriteFile(ctx, path, bytes.NewReader(content), domainfs.WriteOptions{
		Overwrite: force,
		IfMatch:   ifMatch,
	})
	if err != nil {
		if !errors.Is(err, infrafile.ErrPrecondition) {
			c.logger.Error("failed to upload updated file",
				infralogging.String("path", path),
				infralogging.Error(err))
		}
		return err
	}
	return nil
}
