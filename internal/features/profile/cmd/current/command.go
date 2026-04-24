package current

import (
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/spf13/cobra"
)

// CreateCurrentCmd constructs and returns the cobra.Command for the profile current operation.
func CreateCurrentCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("profile-current")
	handler := NewCommand(container.Profile(), l)

	cmd := &cobra.Command{
		Use:   "current",
		Short: "Show the active profile",
		Long:  "Display the name of the currently active profile.",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

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
