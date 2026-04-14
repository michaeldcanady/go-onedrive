package mkdir

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
	"github.com/spf13/cobra"
)

// CreateMkdirCmd constructs and returns the cobra.Command for the drive mkdir operation.
func CreateMkdirCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:               "mkdir <path>",
		Short:             "Create a new directory",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cli.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Path = args[0]
			opts.Stdout = cmd.OutOrStdout()

			// Resolve URI using the factory
			uri, err := container.URIFactory().FromString(opts.Path)
			if err != nil {
				return err
			}
			opts.URI = uri

			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			l, _ := container.Logger().CreateLogger("drive-mkdir")
			handler := NewHandler(container.FS(), l)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	return cmd
}
