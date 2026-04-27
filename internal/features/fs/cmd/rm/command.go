package rm

import (
	cli "github.com/michaeldcanady/go-onedrive/internal/core/cli"
	"github.com/michaeldcanady/go-onedrive/internal/core/di"

	"github.com/spf13/cobra"
)

// CreateRmCmd constructs and returns the cobra.Command for the drive rm operation.
func CreateRmCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("rm")
	handler := NewCommand(container.FS(), container.URIFactory(), l)

	cmd := &cobra.Command{
		Use:               "rm <path>",
		Short:             "Remove files and directories",
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

	cmd.ValidArgsFunction = cli.ProviderPathCompletion(container)

	return cmd
}
