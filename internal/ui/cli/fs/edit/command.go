package edit

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	cli "github.com/michaeldcanady/go-onedrive/internal/ui/cli/fs"

	"github.com/spf13/cobra"
)

// CreateEditCmd constructs and returns the cobra.Command for the edit operation.
func CreateEditCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("edit")
	handler := NewCommand(container.FS(), container.URIFactory(), container.Editor(), l)

	cmd := &cobra.Command{
		Use:               "edit <path>",
		Short:             "Edit a file",
		Long:              "Open a file in your default editor. Changes are synced back when the editor closes.",
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

	cmd.Flags().StringVar(&opts.Editor, "editor", "", "Editor to use (overrides config and environment)")

	return cmd
}
