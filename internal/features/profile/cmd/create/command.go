package create

import (
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/spf13/cobra"
)

// CreateCreateCmd constructs and returns the cobra.Command for the profile create operation.
func CreateCreateCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("profile-create")
	handler := NewCommand(container.Profile(), l)

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new profile",
		Long:  "Create a new profile with the specified name.",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
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
