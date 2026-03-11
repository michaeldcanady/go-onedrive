// Package login provides the command-line interface for authenticating with OneDrive.
package login

import (
	"context"
	"time"

	domainauth "github.com/michaeldcanady/go-onedrive/internal/auth/domain"
	"github.com/michaeldcanady/go-onedrive/internal/interface/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
)

const (
	// FilesReadWriteAllScope is the Microsoft Graph scope for reading and writing all files.
	FilesReadWriteAllScope = "Files.ReadWrite.All"
	// UserReadScope is the Microsoft Graph scope for reading basic user profile information.
	UserReadScope = "User.Read"
	// OfflineAccessScope is the Microsoft Graph scope for requesting refresh tokens.
	OfflineAccessScope = "offline_access"
)

// LoginCmd handles the execution logic for the 'auth login' command.
type LoginCmd struct {
	util.BaseCommand
}

// NewLoginCmd creates a new LoginCmd instance with the provided dependency container.
func NewLoginCmd(container didomain.Container) *LoginCmd {
	return &LoginCmd{
		BaseCommand: util.NewBaseCommand(container, commandName),
	}
}

// Run executes the login flow. It retrieves the current profile, initiates
// authentication via the auth service, and optionally displays the access token.
// It uses specific domain services to decouple from the full container.
func (c *LoginCmd) Run(ctx context.Context, opts Options) error {
	start := time.Now()

	if err := c.Initialize(loggerID); err != nil {
		return util.NewCommandError(c.Name, "failed to initialize command", err)
	}

	c.Log.Info("starting login flow")

	stateSvc := c.Container.State()
	if stateSvc == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "state service is nil")
	}

	profileName, err := stateSvc.Get(domainstate.KeyProfile)
	if err != nil {
		c.Log.Error("failed to retrieve current profile", domainlogger.Error(err))
		return util.NewCommandError(c.Name, "failed to retrieve current profile", err)
	}

	c.Log.Info("resolved current profile",
		domainlogger.String("profile", profileName),
		domainlogger.Bool("force", opts.Force),
		domainlogger.Bool("showToken", opts.ShowToken),
	)

	authService := c.Container.Auth()
	if authService == nil {
		return util.NewCommandErrorWithNameWithMessage(c.Name, "auth service is nil")
	}

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
