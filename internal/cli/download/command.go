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
		Long: `
Download a OneDrive file to a local path.
This is an alias for 'cp onedrive:[src] local:[dst]'.
`,
		Example: `
  # Download a file from OneDrive root
  odc download /notes.txt ./notes.txt

  # Download and overwrite an existing local file
  odc download --force /Documents/report.pdf ./report.pdf
`,

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
