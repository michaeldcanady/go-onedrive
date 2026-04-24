package logout

import (
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/spf13/cobra"
)

// CreateLogoutCmd constructs and returns the cobra.Command for the logout operation.
func CreateLogoutCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("auth-logout")
	handler := NewCommand(
		container.Config(),
		container.Identity(),
		l,
	)

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Sign out from OneDrive",
		Long:  `Sign out from OneDrive for the active profile by clearing the cached authentication tokens.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			c = NewCommandContext(cmd.Context(), &opts)

			return handler.Validate(c)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := handler.Execute(c); err != nil {
				return err
			}
			return handler.Finalize(c)
		},
	}

	cmd.Flags().StringVar(&opts.IdentityID, "id", "", "The specific account to logout (optional)")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Clear all cached credentials for the profile")

	return cmd
}
