package mkdir

import (
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	cli "github.com/michaeldcanady/go-onedrive/internal/ui/cli/fs"

	"github.com/spf13/cobra"
)

// CreateMkdirCmd constructs and returns the cobra.Command for the drive mkdir operation.
func CreateMkdirCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("drive-mkdir")
	handler := NewCommand(container.FS(), container.URIFactory(), l)

	cmd := &cobra.Command{
		Use:               "mkdir <path>",
		Short:             "Create a new directory",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cli.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Path = args[0]
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

	return cmd
}
