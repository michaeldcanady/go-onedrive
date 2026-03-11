package upload

import (
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/cli/cp"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	"github.com/spf13/cobra"
)

const (
	loggerID    = "cli"
	commandName = "upload"

	overwriteFlagName  = "force"
	overwriteFlagShort = "f"
)

// CreateUploadCmd constructs and returns the cobra.Command for the upload operation.
// It initializes flags and sets up the execution logic using CpCmd.
func CreateUploadCmd(c didomain.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:        fmt.Sprintf("%s [src] [dst]", commandName),
		Short:      "Upload a local file to a OneDrive path.",
		Deprecated: "use 'cp local:[src] onedrive:[dst]' instead",
		Long: `
Upload a local file to a OneDrive path.
This is an alias for 'cp local:[src] onedrive:[dst]'.
`,
		Example: `
  # Upload a file to the root of OneDrive
  odc upload ./notes.txt /notes.txt

  # Upload into a folder (basename is appended automatically)
  odc upload ./photo.jpg /Pictures/

  # Upload and overwrite an existing file
  odc upload --force ./report.pdf /Documents/report.pdf
`,

		Args: cobra.ExactArgs(2),

		PreRunE: func(_ *cobra.Command, args []string) error {
			opts.Source = args[0]
			opts.Destination = args[1]
			return opts.Validate()
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			src := opts.Source
			if !strings.HasPrefix(src, "local:") {
				src = "local:" + src
			}

			dst := opts.Destination
			if !strings.HasPrefix(dst, "onedrive:") {
				dst = "onedrive:" + dst
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

	cmd.Flags().BoolVarP(&opts.Overwrite, overwriteFlagName, overwriteFlagShort, false, "Overwrite an existing file at the destination")

	return cmd
}
