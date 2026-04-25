package upload

import (
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	cli "github.com/michaeldcanady/go-onedrive/internal/core/cli"

	"github.com/spf13/cobra"
)

// CreateUploadCmd constructs and returns the cobra.Command for the drive upload operation.
func CreateUploadCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("upload")
	handler := NewCommand(container.FS(), container.URIFactory(), l)

	cmd := &cobra.Command{
		Use:               "upload <source> <destination>",
		Short:             "Upload files and directories",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: cli.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Source = args[0]
			opts.Destination = args[1]
			opts.Stdout = cmd.OutOrStdout()

			c = &CommandContext{
				Ctx:     cmd.Context(),
				Options: opts,
			}

			return handler.Validate(c)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := handler.Execute(c); err != nil {
				return err
			}
			return handler.Finalize(c)
		},
	}

	cmd.ValidArgsFunction = cli.ProviderPathCompletion(container)
	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "upload directories recursively")

	return cmd
}
