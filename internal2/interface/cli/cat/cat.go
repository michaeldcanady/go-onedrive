package cat

import (
	"context"
	"fmt"
	"io"
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
	commandName = "cat"
)

func CreateCatCmd(c di.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s [path]", commandName),
		Short: "List items in a OneDrive path",
		Args:  cobra.MaximumNArgs(1),

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

			logger.Info("starting cat command")

			// Filesystem service
			fsSvc := c.FS()
			if fsSvc == nil {
				logger.Error("filesystem service is nil")
				return util.NewCommandErrorWithNameWithMessage(commandName, "filesystem service is nil")
			}

			// Resolve path
			path := ""
			if len(args) > 0 {
				path = args[0]
			}

			logger.Debug("path resolved", infralogging.String("path", path))
			if strings.TrimSpace(path) == "" {
				logger.Error("path is empty")
				return util.NewCommandErrorWithNameWithMessage(commandName, "path is empty")
			}

			reader, err := fsSvc.ReadFile(ctx, path, fs.ReadOptions{})
			if err != nil {
				logger.Error("failed to read file", infralogging.Error(err))
				return util.NewCommandErrorWithNameWithMessage(commandName, "unable to read path contents")
			}
			defer reader.Close()

			_, err = io.Copy(cmd.OutOrStdout(), reader)
			if err != nil {
				logger.Error("failed to write file contents", infralogging.Error(err))
				return util.NewCommandErrorWithNameWithMessage(commandName, "failed to write file contents")
			}

			logger.Info("cat command completed",
				infralogging.Duration("duration", time.Since(start)),
			)

			return nil
		},
	}

	return cmd
}
