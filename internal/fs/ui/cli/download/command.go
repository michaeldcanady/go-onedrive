package download

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
	"github.com/spf13/cobra"
)

// CreateDownloadCmd constructs and returns the cobra.Command for the drive download operation.
func CreateDownloadCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:               "download <remote_path> <local_path>",
		Short:             "Download files and directories from OneDrive",
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
			l, _ := container.Logger().CreateLogger("drive-download")
			handler := NewHandler(container.FS(), l)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "download directories recursively")

	return cmd
}
