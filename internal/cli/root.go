package cli

import "github.com/spf13/cobra"

const (
	rootCmdName   = "onedrive"
	rootCmdShort  = "OneDrive CLI"
	authCmdName   = "auth"
	authCmdShort  = "Authenticate with OneDrive"
	loginCmdName  = "login"
	loginCmdShort = "Login to OneDrive"
)

var rootCmd = &cobra.Command{
	Use:   rootCmdName,
	Short: rootCmdShort,
}
