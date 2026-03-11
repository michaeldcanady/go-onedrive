package upload

import (
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/cp"
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
		Use:        fmt.Sprintf("%s <SOURCE> <DESTINATION>", commandName),
		Short:      "Upload a local file to a OneDrive path.",
		Deprecated: "use 'cp local:<SOURCE> onedrive:<DESTINATION>' instead",
		Long: `You can upload a local file to a specified OneDrive path. This command is a
convenient alias for the 'cp' command using the 'local:' prefix for the
source and 'onedrive:' for the destination.`,
		Example: `  # Upload a file to the OneDrive root
  odc upload ./presentation.pptx /presentation.pptx

  # Upload a file and overwrite any existing file with the same name
  odc upload --force ./data.csv /ProjectX/data.csv`,

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
