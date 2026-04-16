package delete

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/pkg/args"
	"github.com/spf13/cobra"
)

// CreateDeleteCmd constructs and returns the cobra.Command for the profile delete operation.
func CreateDeleteCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("profile-delete")
	handler := NewCommand(container.Profile(), l)

	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a profile",
		Long:  "Permanently delete a profile and its associated configuration and authentication tokens.",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, argsSlice []string) error {
			if err := args.Bind(argsSlice, &opts); err != nil {
				return err
			}
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
