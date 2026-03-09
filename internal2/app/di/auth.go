package di

import (
	"path/filepath"

	appaccount "github.com/michaeldcanady/go-onedrive/internal2/app/account"
	appauth "github.com/michaeldcanady/go-onedrive/internal2/app/auth"
	appcache "github.com/michaeldcanady/go-onedrive/internal2/app/cache"
	appprofile "github.com/michaeldcanady/go-onedrive/internal2/app/profile"
	appstate "github.com/michaeldcanady/go-onedrive/internal2/app/state"

	domainaccount "github.com/michaeldcanady/go-onedrive/internal2/domain/account"
	domainauth "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	domainprofile "github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
	domainstate "github.com/michaeldcanady/go-onedrive/internal2/domain/state"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/auth/msal"
	infraprofile "github.com/michaeldcanady/go-onedrive/internal2/infra/profile"
	infrastate "github.com/michaeldcanady/go-onedrive/internal2/infra/state"
)

// Auth implements [di.Container].
func (c *Container) Auth() domainauth.AuthService {
	c.authOnce.Do(func() {
		c.authService = c.newAuthService()
	})

	return c.authService
}

func (c *Container) newAuthService() domainauth.AuthService {
	credentialFactory := msal.NewMSALCredentialFactory()

	return appauth.NewService2(c.authCache(), c.Config(), c.State(), c.getLogger("auth"), credentialFactory, c.Account())
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

// Profile implements [di.Container].
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
