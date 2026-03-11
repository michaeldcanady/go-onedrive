package login

import (
	"context"
	"fmt"
	"time"

	domainauth "github.com/michaeldcanady/go-onedrive/internal/auth/domain"
	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
)

const (
	FilesReadWriteAllScope = "Files.ReadWrite.All"
	UserReadScope          = "User.Read"
	OfflineAccessScope     = "offline_access"
)

type LoginCmd struct {
	util.BaseCommand
}

func NewLoginCmd(container didomain.Container) *LoginCmd {
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
		c.Log.Error("failed to get current profile", domainlogger.Error(err))
		return util.NewCommandErrorWithNameWithError(c.Name, err)
	}

	c.Log.Info("resolved current profile",
		domainlogger.String("profile", profileName),
		domainlogger.Bool("force", opts.Force),
		domainlogger.Bool("showToken", opts.ShowToken),
	)

	authService := c.Container.Auth()

	loginOpts := domainauth.LoginOptions{
		Force: opts.Force,
		Scopes: []string{
			FilesReadWriteAllScope,
			UserReadScope,
		},
		EnableCAE: true,
	}

	c.Log.Info("initiating authentication",
		domainlogger.String("profile", profileName),
		domainlogger.Bool("force", loginOpts.Force),
		domainlogger.Strings("scopes", loginOpts.Scopes),
		domainlogger.Bool("enableCAE", loginOpts.EnableCAE),
	)

	result, err := authService.Login(ctx, profileName, loginOpts)
	if err != nil {
		c.Log.Error("authentication failed",
			domainlogger.String("profile", profileName),
			domainlogger.Error(err),
		)
		fmt.Fprintf(opts.Stderr, "Login failed for profile %q: %v\n", profileName, err)
		return util.NewCommandError(c.Name, "failed authentication", err)
	}

	c.Log.Info("authentication successful",
		domainlogger.String("profile", profileName),
		domainlogger.Bool("tokenDisplayed", opts.ShowToken),
	)

	fmt.Fprintf(opts.Stdout, "Successfully logged into profile %q\n", profileName)

	if opts.ShowToken {
		fmt.Fprintf(opts.Stdout, "Access Token: %s\n", result.AccessToken)
	}

	c.Log.Info("login flow completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
