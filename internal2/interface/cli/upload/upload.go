package upload

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/spf13/cobra"
)

const (
	loggerID    = "cli"
	commandName = "upload"

	overwriteFlagName  = "force"
	overwriteFlagShort = "f"
	overwriteFlagUsage = ""
)

func CreateUploadCmd(c di.Container) *cobra.Command {
	var (
		overwrite bool
	)

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [src] [dst]", commandName),
		Short: "Upload a local file to a OneDrive path.",
		Long: `
Upload a local file to a OneDrive path.

This command copies a file from your local filesystem into your OneDrive
account. It behaves similarly to the Unix 'cp' command, but with cloud-aware
semantics:

  • The first argument (src) must be a path to a local file.
  • The second argument (dst) is the destination path in OneDrive.
  • If the destination ends with a slash ("/"), the source file's basename
    is automatically appended.
  • Existing files at the destination path are overwritten unless OneDrive
    prevents it.
  • Parent folders must already exist unless your OneDrive configuration
    supports implicit folder creation.

This command does not currently support uploading directories, recursive
uploads, or glob patterns. For multi-file uploads, run this command in a loop
or use a higher-level automation script.

Authentication:
You must be logged in (via 'onedrive auth login') before using this command.
`,
		Example: `
  # Upload a file to the root of OneDrive
  onedrive upload ./notes.txt /notes.txt

  # Upload into a folder (basename is appended automatically)
  onedrive upload ./photo.jpg /Pictures/

  # Upload and overwrite an existing file
  onedrive upload ./report.pdf /Documents/report.pdf

  # Upload using a relative OneDrive path
  onedrive upload ./todo.md Documents/todo.md

  # Upload a file whose name should be preserved
  onedrive upload ./archive.tar.gz /Backups/archive.tar.gz
`,

		Args: cobra.ExactArgs(2),

		RunE: func(cmd *cobra.Command, args []string) error {
			start := time.Now()

			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			logger, err := util.EnsureLogger(c, loggerID)
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger = logger.WithContext(ctx).With(logging.String("correlationID", util.CorrelationIDFromContext(ctx)))

			fsSvc := c.FS()
			if fsSvc == nil {
				return util.NewCommandErrorWithNameWithMessage(commandName, "filesystem service is nil")
			}

			src := args[0]
			dst := args[1]

			if strings.HasSuffix(dst, string(os.PathSeparator)) || strings.HasSuffix(dst, "/") {
				name := filepath.Base(src)
				dst = fmt.Sprintf("%s%s", dst, name)
			}

			file, err := os.OpenFile(src, os.O_RDONLY, 0)
			if err != nil {
				return util.NewCommandError(commandName, "failed to open file", err)
			}
			defer file.Close()

			_, err = fsSvc.WriteFile(ctx, dst, file, fs.WriteOptions{Overwrite: overwrite})
			if err != nil {
				return util.NewCommandError(commandName, "failed to upload file", err)
			}

			logger.Info("file updated successfully",
				infralogging.Duration("duration", time.Since(start)),
			)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&overwrite, overwriteFlagName, overwriteFlagShort, false, overwriteFlagUsage)

	return cmd
}
