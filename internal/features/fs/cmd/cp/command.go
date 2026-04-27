package cp

import (
	cli "github.com/michaeldcanady/go-onedrive/internal/core/cli"
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/spf13/cobra"
)

// CreateCpCmd constructs and returns the cobra.Command for the drive cp operation.
func CreateCpCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("cp")
	handler := NewCommand(container.FS(), container.URIFactory(), l)

	cmd := &cobra.Command{
		Use:               "cp <source> <destination>",
		Short:             "Copy files and directories",
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

	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "copy directories recursively")

	return cmd
}
