package use

import (
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/spf13/cobra"
)

// CreateUseCmd constructs and returns the cobra.Command for the profile use operation.
func CreateUseCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("profile-use")
	handler := NewCommand(container.Profile(), l)

	cmd := &cobra.Command{
		Use:   "use <name>",
		Short: "Set the active profile",
		Long:  "Specify a profile name to be used for subsequent commands.",
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
