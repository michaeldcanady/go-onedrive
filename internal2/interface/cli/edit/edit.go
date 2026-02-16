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
	"github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
	"github.com/spf13/cobra"
)

const (
	loggerID    = "cli"
	commandName = "edit"
)

func CreateEditCmd(c di.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [path]", commandName),
		Short: "Edit a OneDrive file in your local editor",
		Args:  cobra.ExactArgs(1),

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

			path := args[0]

			reader, err := fsSvc.ReadFile(ctx, path, fs.ReadOptions{})
			if err != nil {
				return util.NewCommandError(commandName, "failed to read file", err)
			}
			defer reader.Close()

			origHash := sha256.New()

			name := Name(path)
			ext := filepath.Ext(path)

			editorSvc := NewEditorService(c.EnvironmentService(), logger).
				WithIO(cmd.InOrStdin(), cmd.OutOrStdout(), cmd.ErrOrStderr())

			editedBytes, tmpPath, err := editorSvc.LaunchTempFile(fmt.Sprintf("%s-edit-", name), ext, io.TeeReader(reader, origHash))
			defer os.Remove(tmpPath)
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			origHashSum := hex.EncodeToString(origHash.Sum(nil))

			editedHash := sha256.Sum256(editedBytes)

			if origHashSum == hex.EncodeToString(editedHash[:]) {
				logger.Info("no changes detected")
				return nil
			}

			_, err = fsSvc.WriteFile(ctx, path, bytes.NewReader(editedBytes), fs.WriteOptions{})
			if err != nil {
				return util.NewCommandError(commandName, "failed to write updated file", err)
			}

			logger.Info("file updated successfully",
				infralogging.Duration("duration", time.Since(start)),
			)

			return nil
		},
	}

	return cmd
}
