package logout

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

// CreateLogoutCmd constructs and returns the cobra.Command for the logout operation.
func CreateLogoutCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Sign out from your OneDrive profile",
		Long: `This command clears the authentication state for your active profile.
Any cached tokens will be invalidated, and subsequent operations will require
re-authentication.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			l, _ := container.Logger().CreateLogger("auth-logout")
			handler := NewHandler(
				container.Config(),
				container.Identity(),
				l,
			)

			return handler.Handle(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Force clear all cached credentials")

	return cmd
}
