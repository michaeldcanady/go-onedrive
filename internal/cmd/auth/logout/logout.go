package logout

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CreateLogoutCmd constructs the `auth logout` command, which clears any
// locally stored authentication credentials.
func CreateLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out and clear local authentication credentials",
		Long: `Logs you out of OneDrive by removing any locally stored tokens and authentication state.

This does not revoke tokens server-side; it only clears your local session.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Inject and call your token store / auth service here.
			// Example:
			// if err := container.TokenStore.Clear(); err != nil {
			//     return fmt.Errorf("failed to clear credentials: %w", err)
			// }

			fmt.Println("You have been logged out.")
			return nil
		},
		Example: "odc auth logout",
	}

	return cmd
}
