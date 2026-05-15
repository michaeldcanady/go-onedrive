package edit

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/michaeldcanady/go-onedrive/internal/features/vfs"
)

// Validate performs initial validation of the command options.
func (c *Command) Validate(ctx *CommandContext) error {
	if ctx.Options.Path == "" {
		return fmt.Errorf("path is required")
	}
	return nil
}

// Resolve performs argument resolution.
func (c *Command) Resolve(ctx *CommandContext) error {
	return c.BaseResolve(ctx)
}

// Execute performs the core business logic of the command.
func (c *Command) Execute(ctx *CommandContext) error {
	// 1. Get initial remote state
	node, err := c.fS.Stat(ctx.Ctx, ctx.Options.Path)
	if err != nil {
		return fmt.Errorf("failed to stat remote file: %w", err)
	}

	// 2. Download to temp
	reader, err := c.fS.Read(ctx.Ctx, ctx.Options.Path)
	if err != nil {
		return fmt.Errorf("failed to read remote file: %w", err)
	}
	defer reader.Close()

	tempDir, err := os.MkdirTemp("", "odc-edit-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	tempPath := filepath.Join(tempDir, filepath.Base(ctx.Options.Path))
	f, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, reader); err != nil {
		f.Close()
		return err
	}
	f.Close()

	// 3. Initial stat of local file
	initialInfo, err := os.Stat(tempPath)
	if err != nil {
		return err
	}

	// 4. Launch editor
	if err := c.editor.Open(ctx.Ctx, tempPath); err != nil {
		return fmt.Errorf("editor session failed: %w", err)
	}

	// 5. Check for modification
	finalInfo, err := os.Stat(tempPath)
	if err != nil {
		return err
	}

	if finalInfo.ModTime().After(initialInfo.ModTime()) {
		fmt.Println("Changes detected, uploading...")
		f, err := os.Open(tempPath)
		if err != nil {
			return err
		}
		defer f.Close()

		var writeOpts []vfs.WriteOption
		if !ctx.Options.Force && node.ETag != "" {
			writeOpts = append(writeOpts, vfs.WithIfMatch(node.ETag))
		}

		if err := c.fS.Write(ctx.Ctx, ctx.Options.Path, f, writeOpts...); err != nil {
			return fmt.Errorf("failed to write changes (use --force to overwrite upstream changes): %w", err)
		}
		fmt.Println("Changes synced successfully")
	} else {
		fmt.Println("No changes detected")
	}

	return nil
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
