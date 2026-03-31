package download

import (
	"github.com/michaeldcanady/go-onedrive/internal/feature/di"
	"github.com/spf13/cobra"
)

// CreateDownloadCmd constructs and returns the cobra.Command for the drive download operation.
func CreateDownloadCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "download <remote_path> <local_path>",
		Short: "Download files and directories from OneDrive",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Source = args[0]
			opts.Destination = args[1]
			opts.Stdout = cmd.OutOrStdout()

			if err := opts.Validate(); err != nil {
				return err
			}

			l, _ := container.Logger().CreateLogger("drive-download")
			handler := NewHandler(container.FS(), l)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "download directories recursively")

	return cmd
}
