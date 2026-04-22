package logout

import (
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	"github.com/spf13/cobra"
)

// CreateLogoutCmd constructs and returns the cobra.Command for the logout operation.
func CreateLogoutCmd(container di.Container) *cobra.Command {
	var opts Options

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

			return handler.Validate(cmd.Context(), &opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return handler.Execute(cmd.Context(), opts)
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return handler.Finalize(cmd.Context(), opts)
		},
	}

	return cmd
}
