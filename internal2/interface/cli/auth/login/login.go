package login

import (
	"context"
	"fmt"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	logger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
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

type LoginCmd struct {
	util.BaseCommand
}

func NewLoginCmd(container di.Container) *LoginCmd {
	return &LoginCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *LoginCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return err
	}

	c.Log.Info("starting login flow")

	profileName, err := c.Container.State().GetCurrentProfile()
	if err != nil {
		c.Log.Error("failed to get current profile", logger.Error(err))
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	c.Log.Info("resolved current profile",
		logger.String("profile", profileName),
		logger.Bool("force", opts.Force),
		logger.Bool("showToken", opts.ShowToken),
	)

	authService := c.Container.Auth()

	loginOpts := auth.LoginOptions{
		Force: opts.Force,
		Scopes: []string{
			FilesReadWriteAllScope,
			UserReadScope,
		},
		EnableCAE: true,
	}

	c.Log.Info("initiating authentication",
		logger.String("profile", profileName),
		logger.Bool("force", loginOpts.Force),
		logger.Strings("scopes", loginOpts.Scopes),
		logger.Bool("enableCAE", loginOpts.EnableCAE),
	)

	result, err := authService.Login(ctx, profileName, loginOpts)
	if err != nil {
		c.Log.Error("authentication failed",
			logger.String("profile", profileName),
			logger.Error(err),
		)
		fmt.Fprintf(opts.Stderr, "Login failed for profile %q: %v\n", profileName, err)
		return util.NewCommandError(c.Name, "failed authentication", err)
	}

	c.Log.Info("authentication successful",
		logger.String("profile", profileName),
		logger.Bool("tokenDisplayed", opts.ShowToken),
	)

	fmt.Fprintf(opts.Stdout, "Successfully logged into profile %q\n", profileName)

	if opts.ShowToken {
		fmt.Fprintf(opts.Stdout, "Access Token: %s\n", result.AccessToken)
	}

	c.Log.Info("login flow completed successfully",
		logger.Duration("duration", time.Since(start)),
	)
	return nil
}

func CreateLoginCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with OneDrive",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			return NewLoginCmd(container).Run(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVar(&opts.ShowToken, showTokenLongFlag, false, showTokenUsage)
	cmd.Flags().BoolVarP(&opts.Force, forceLongFlag, forceShortFlag, false, forceUsage)

	return cmd
}
