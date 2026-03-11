package download

import (
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/cli/cp"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	loggerID    = "cli"
	commandName = "download"
)

// CreateDownloadCmd constructs and returns the cobra.Command for the download operation.
func CreateDownloadCmd(c didomain.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:        fmt.Sprintf("%s [src] [dst]", commandName),
		Short:      "Download a OneDrive file to a local path.",
		Deprecated: "use 'cp onedrive:[src] local:[dst]' instead",
		Long: `You can download a file from OneDrive to your local machine. This command is
a convenient alias for the 'cp' command using the 'onedrive:' prefix for the
source and 'local:' for the destination.`,
		Example: `  # Download a file from OneDrive to your current directory
  odc download /Shared/report.pdf ./report.pdf

  # Download and overwrite a local file if it already exists
  odc download --force /Photos/vacation.jpg ./vacation.jpg`,

		Args: cobra.ExactArgs(2),

		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Source = args[0]
			opts.Destination = args[1]
			return opts.Validate()
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			src := opts.Source
			if !strings.HasPrefix(src, "onedrive:") {
				src = "onedrive:" + src
			}

			dst := opts.Destination
			if !strings.HasPrefix(dst, "local:") {
				dst = "local:" + dst
			}

			cpOpts := cp.Options{
				Source:    src,
				Dest:      dst,
				Overwrite: opts.Overwrite,
				Stdin:     cmd.InOrStdin(),
				Stdout:    cmd.OutOrStdout(),
				Stderr:    cmd.ErrOrStderr(),
			}

			cpCmd := cp.NewCpCmd(c)
			return cpCmd.Run(cmd.Context(), cpOpts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Overwrite, "force", "f", false, "Overwrite an existing local file")

	return cmd
}
