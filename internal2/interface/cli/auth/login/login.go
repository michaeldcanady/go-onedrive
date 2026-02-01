package login

import (
	"context"

	applogging "github.com/michaeldcanady/go-onedrive/internal2/app/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
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
	loggerID    = "cli"
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

			logger, err := ensureLogger(container)
			if err != nil {
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("starting login flow")

			profileName, err := container.State().GetCurrentProfile()
			if err != nil {
				logger.Error("failed to get current profile", infralogging.Error(err))
				return util.NewCommandErrorWithNameWithError(commandName, err)
			}

			logger.Info("resolved current profile",
				infralogging.String("profile", profileName),
				infralogging.Bool("force", force),
				infralogging.Bool("showToken", showToken),
			)

			authService := container.Auth()

			opts := auth.LoginOptions{
				Force: force,
				Scopes: []string{
					FilesReadWriteAllScope,
					UserReadScope,
				},
				EnableCAE: true,
			}

			logger.Info("initiating authentication",
				infralogging.String("profile", profileName),
				infralogging.Bool("force", opts.Force),
				infralogging.Strings("scopes", opts.Scopes),
				infralogging.Bool("enableCAE", opts.EnableCAE),
			)

			result, err := authService.Login(ctx, profileName, opts)
			if err != nil {
				logger.Error("authentication failed",
					infralogging.String("profile", profileName),
					infralogging.Error(err),
				)
				return util.NewCommandError(commandName, "failed authentication", err)
			}

			logger.Info("authentication successful",
				infralogging.String("profile", profileName),
				infralogging.Bool("tokenDisplayed", showToken),
			)

			if showToken {
				// Only print token if explicitly requested
				cmd.Printf("Access Token: %s\n", result.AccessToken)
			}

			logger.Info("login flow completed successfully")
			return nil
		},
	}

	cmd.Flags().BoolVar(&showToken, showTokenLongFlag, false, showTokenUsage)
	cmd.Flags().BoolVarP(&force, forceLongFlag, forceShortFlag, false, forceUsage)

	return cmd
}

func ensureLogger(c di.Container) (infralogging.Logger, error) {
	logger, err := c.Logger().GetLogger(loggerID)
	if err == applogging.ErrUnknownLogger {
		return c.Logger().CreateLogger(loggerID)
	}
	return logger, err
}
