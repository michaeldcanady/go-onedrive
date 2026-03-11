package login

import (
	"context"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
	"time"

	domainauth "github.com/michaeldcanady/go-onedrive/internal/auth/domain"
	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	domaindi "github.com/michaeldcanady/go-onedrive/internal/di/domain"
)

const (
	FilesReadWriteAllScope = "Files.ReadWrite.All"
	UserReadScope          = "User.Read"
	OfflineAccessScope     = "offline_access"
)

type LoginCmd struct {
	util.BaseCommand
}

func NewLoginCmd(container domaindi.Container) *LoginCmd {
	return &LoginCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

func (c *LoginCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("starting login flow")

	profileName, err := c.Container.State().Get(domainstate.KeyProfile)
	if err != nil {
		return util.NewCommandError(c.Name, "failed to retrieve current profile", err)
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
		return util.NewCommandError(c.Name, "failed authentication", err)
	}

	c.Log.Info("authentication successful",
		domainlogger.String("profile", profileName),
		domainlogger.Bool("tokenDisplayed", opts.ShowToken),
	)

	c.RenderSuccess(opts.Stdout, "Successfully logged into profile %q\n", profileName)

	if opts.ShowToken {
		c.RenderMessage(opts.Stdout, "Access Token: %s\n", result.AccessToken)
	}

	c.Log.Info("login flow completed successfully",
		domainlogger.Duration("duration", time.Since(start)),
	)
	return nil
}
