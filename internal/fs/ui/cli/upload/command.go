package upload

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
	"github.com/spf13/cobra"
)

// CreateUploadCmd constructs and returns the cobra.Command for the drive upload operation.
func CreateUploadCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:               "upload <local_path> <remote_path>",
		Short:             "Upload files and directories to OneDrive",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: cli.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Source = args[0]
			opts.Stdout = cmd.OutOrStdout()

			// Resolve URI using the factory
			sourceURI, err := container.URIFactory().FromString(opts.Source)
			if err != nil {
				return err
			}
			opts.SourceURI = sourceURI

			opts.Destination = args[1]
			// Resolve URI using the factory
			destinationURI, err := container.URIFactory().FromString(opts.Destination)
			if err != nil {
				return err
			}
			opts.DestinationURI = destinationURI

			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			l, _ := container.Logger().CreateLogger("drive-upload")
			handler := NewHandler(container.FS(), l)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "upload directories recursively")

	return cmd
}
