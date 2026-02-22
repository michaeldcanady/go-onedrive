package upload

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

// UploadCmd handles the execution logic for the 'upload' command.
// It coordinates reading a local file and writing it to the OneDrive service.
type UploadCmd struct {
	container di.Container
	logger    infralogging.Logger
}

// NewUploadCmd creates a new UploadCmd instance with the provided dependency container.
func NewUploadCmd(container di.Container) *UploadCmd {
	return &UploadCmd{
		container: container,
	}
}

// WithLogger allows injecting a logger into UploadCmd for testing.
func (c *UploadCmd) WithLogger(logger infralogging.Logger) *UploadCmd {
	c.logger = logger
	return c
}

// Run executes the upload lifecycle.
// 1. Ensures a logger is available.
// 2. Resolves the final destination path.
// 3. Opens the local source file.
// 4. Uploads the file via the filesystem service.
func (c *UploadCmd) Run(ctx context.Context, opts Options) error {
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

	c.logger.Info("starting upload command",
		infralogging.String("source", opts.Source),
		infralogging.String("destination", opts.Destination),
		infralogging.Bool("overwrite", opts.Overwrite),
	)

	fsSvc := c.container.FS()
	if fsSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(commandName, "filesystem service is nil")
	}

	// 1. Resolve Destination
	dst := c.resolveDestination(opts.Source, opts.Destination)

	// 2. Open Local File
	file, err := os.OpenFile(opts.Source, os.O_RDONLY, 0)
	if err != nil {
		return util.NewCommandError(commandName, "failed to open local file", err)
	}
	defer file.Close()

	// 3. Perform Upload
	c.logger.Debug("initiating file write", infralogging.String("resolvedPath", dst))
	_, err = fsSvc.WriteFile(ctx, dst, file, domainfs.WriteOptions{Overwrite: opts.Overwrite})
	if err != nil {
		return util.NewCommandError(commandName, "failed to upload file", err)
	}

	c.logger.Info("upload completed successfully",
		infralogging.String("path", dst),
		infralogging.Duration("duration", time.Since(start)),
	)

	return nil
}

// resolveDestination appends the source filename to the destination if the destination
// indicates it's a directory (ends with a slash).
func (c *UploadCmd) resolveDestination(src, dst string) string {
	if strings.HasSuffix(dst, string(os.PathSeparator)) || strings.HasSuffix(dst, "/") {
		name := filepath.Base(src)
		// Ensure we don't have double slashes if possible, though OneDrive API usually handles it.
		if !strings.HasSuffix(dst, "/") && !strings.HasSuffix(dst, string(os.PathSeparator)) {
			dst += "/"
		}
		dst += name
	}
	return dst
}
