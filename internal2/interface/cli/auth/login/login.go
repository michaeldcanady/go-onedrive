package login

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	"github.com/spf13/cobra"
)

const (
	FilesReadWriteAllScope = "Files.ReadWrite.All"
	UserReadScope          = "User.Read"
	OfflineAccessScope     = "offline_access"

	showTokenLongFlag = "show-token"
	showTokenUsage    = "Display the access token after login"

	forceLongFlag  = "force"
	forceShortFlag = "f"
	forceUsage     = "Force re-authentication even if a valid profile exists"

	commandName = "login"
)

func CreateLoginCmd(container di.Container) *cobra.Command {
	var (
		showToken bool
		force     bool
	)

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with OneDrive",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			authService := container.Auth()

			opts := auth.LoginOptions{
				Force: force,
				Scopes: []string{
					FilesReadWriteAllScope,
					UserReadScope,
				},
				EnableCAE: true,
			}

			result, err := authService.Login(ctx, "default", opts)
			if err != nil {
				return NewCommandError(commandName, "failed authentication", err)
			}

			if showToken {
				fmt.Printf("Access Token: %s \n", result.AccessToken)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&showToken, showTokenLongFlag, false, showTokenUsage)
	cmd.Flags().BoolVarP(&force, forceLongFlag, forceShortFlag, false, forceUsage)

	return cmd
}
