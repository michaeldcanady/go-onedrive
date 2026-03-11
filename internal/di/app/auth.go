package app

import (
	"path/filepath"

	appaccount "github.com/michaeldcanady/go-onedrive/internal/account/app"
	authapp "github.com/michaeldcanady/go-onedrive/internal/auth/app"
	infraauth "github.com/michaeldcanady/go-onedrive/internal/auth/infra"
	appcache "github.com/michaeldcanady/go-onedrive/internal/cache/app"
	appprofile "github.com/michaeldcanady/go-onedrive/internal/profile/app"
	appstate "github.com/michaeldcanady/go-onedrive/internal/state/app"

	domainaccount "github.com/michaeldcanady/go-onedrive/internal/account/domain"
	domainauth "github.com/michaeldcanady/go-onedrive/internal/auth/domain"
	domainprofile "github.com/michaeldcanady/go-onedrive/internal/profile/domain"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"

	infraprofile "github.com/michaeldcanady/go-onedrive/internal/profile/infra"
	infrastate "github.com/michaeldcanady/go-onedrive/internal/state/infra"
)

// Auth implements [didomain.Container].
func (c *Container) Auth() domainauth.AuthService {
	c.authOnce.Do(func() {
		c.authService = c.newAuthService()
	})

	return c.authService
}

func (c *Container) newAuthService() domainauth.AuthService {
	credentialFactory := infraauth.NewMSALCredentialFactory()

	return authapp.NewService(c.authCache(), c.Config(), c.State(), c.getLogger("auth"), credentialFactory, c.Account())
}

func (c *Container) Account() domainaccount.Service {
	c.accountOnce.Do(func() {
		c.accountService = c.newAccountService()
	})
	return c.accountService
}

func (c *Container) newAccountService() domainaccount.Service {
	return appaccount.New(c.accountCache(), c.getLogger("account"))
}

// Profile implements [didomain.Container].
func (c *Container) Profile() domainprofile.ProfileService {
	c.profileOnce.Do(func() {
		c.profileService = c.newProfileService()
	})

	return c.profileService
}

func (c *Container) newProfileService() domainprofile.ProfileService {
	env := c.EnvironmentService()

	// ~/.config/odc
	profileBaseDir, err := env.ConfigDir()
	if err != nil {
		panic(err)
	}

	// Infra repository
	repo := infraprofile.NewFSProfileService(profileBaseDir)

	// App service (repo only)
	return appprofile.New(
		c.getLogger("profile"),
		repo,
	)
}

func (c *Container) State() domainstate.Service {
	c.stateOnce.Do(func() {
		c.stateService = c.newStateService()
	})
	return c.stateService
}

func (c *Container) newStateService() domainstate.Service {
	env := c.EnvironmentService()
	stateDir, _ := env.StateDir()
	statePath := filepath.Join(stateDir, stateFileName)

	serializer := &appcache.JSONSerializerDeserializer[domainstate.State]{}
	repo := infrastate.NewRepository(statePath, serializer)

	return appstate.NewService(repo)
}
