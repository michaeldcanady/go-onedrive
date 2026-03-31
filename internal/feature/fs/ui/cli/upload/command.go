package upload

import (
	"github.com/michaeldcanady/go-onedrive/internal/feature/di"
	"github.com/spf13/cobra"
)

// CreateUploadCmd constructs and returns the cobra.Command for the drive upload operation.
func CreateUploadCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "upload <local_path> <remote_path>",
		Short: "Upload files and directories to OneDrive",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Source = args[0]
			opts.Destination = args[1]
			opts.Stdout = cmd.OutOrStdout()

			if err := opts.Validate(); err != nil {
				return err
			}

			l, _ := container.Logger().CreateLogger("drive-upload")
			handler := NewHandler(container.FS(), l)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "upload directories recursively")

	return cmd
}
