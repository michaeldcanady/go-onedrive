package upload

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/spf13/cobra"
)

const (
	loggerID    = "cli"
	commandName = "upload"

	overwriteFlagName  = "force"
	overwriteFlagShort = "f"
)

// CreateUploadCmd constructs and returns the cobra.Command for the upload operation.
// It initializes flags and sets up the execution logic using UploadCmd.
func CreateUploadCmd(c di.Container) *cobra.Command {
	var opts Options

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
  • Existing files at the destination path are overwritten if '--force' is used.
  • Parent folders must already exist.

Authentication:
You must be logged in (via 'onedrive auth login') before using this command.
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
			uploadCmd := NewUploadCmd(c)
			return uploadCmd.Run(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Overwrite, overwriteFlagName, overwriteFlagShort, false, "Overwrite an existing file at the destination")

	return cmd
}
