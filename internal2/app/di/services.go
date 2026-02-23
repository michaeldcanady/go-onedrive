package di

import (
	"path/filepath"

	appaccount "github.com/michaeldcanady/go-onedrive/internal2/app/account"
	appauth "github.com/michaeldcanady/go-onedrive/internal2/app/auth"
	appcache "github.com/michaeldcanady/go-onedrive/internal2/app/cache"
	appconfig "github.com/michaeldcanady/go-onedrive/internal2/app/config"
	appdrive "github.com/michaeldcanady/go-onedrive/internal2/app/drive"
	appeditor "github.com/michaeldcanady/go-onedrive/internal2/app/editor"
	appfs "github.com/michaeldcanady/go-onedrive/internal2/app/fs"
	applogging "github.com/michaeldcanady/go-onedrive/internal2/app/common/logging"
	appprofile "github.com/michaeldcanady/go-onedrive/internal2/app/profile"
	appstate "github.com/michaeldcanady/go-onedrive/internal2/app/state"
	"github.com/michaeldcanady/go-onedrive/internal2/app/common/environment"

	domainaccount "github.com/michaeldcanady/go-onedrive/internal2/domain/account"
	domainauth "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	domainenv "github.com/michaeldcanady/go-onedrive/internal2/domain/common/environment"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	domainconfig "github.com/michaeldcanady/go-onedrive/internal2/domain/config"
	domaindrive "github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	domaineditor "github.com/michaeldcanady/go-onedrive/internal2/domain/editor"
	domainfile "github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	domainprofile "github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
	domainstate "github.com/michaeldcanady/go-onedrive/internal2/domain/state"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/auth/msal"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	infraconfig "github.com/michaeldcanady/go-onedrive/internal2/infra/config"
	infraeditor "github.com/michaeldcanady/go-onedrive/internal2/infra/editor"
	infrafile "github.com/michaeldcanady/go-onedrive/internal2/infra/file"
	infraprofile "github.com/michaeldcanady/go-onedrive/internal2/infra/profile"
	infrastate "github.com/michaeldcanady/go-onedrive/internal2/infra/state"
)

// Drive implements [di.Container].
func (c *Container) Drive() domaindrive.DriveService {
	c.driveOnce.Do(func() {
		c.driveService = c.newDriveService()
	})
	return c.driveService
}

func (c *Container) newDriveService() domaindrive.DriveService {
	loggerService := c.Logger()
	logger, err := loggerService.CreateLogger("drive")
	if err != nil {
		logger = infralogging.NewNoopLogger()
	}

	return appdrive.NewDriveService(c.clientProvider(), logger)
}

// File implements [di.Container].
func (c *Container) File() domainfile.FileService {
	c.fileOnce.Do(func() {
		c.fileService = c.newFileService()
	})
	return c.fileService
}

func (c *Container) newFileService() domainfile.FileService {
	loggerService := c.Logger()
	logger, err := loggerService.CreateLogger("file")
	if err != nil {
		logger = infralogging.NewNoopLogger()
	}

	return infrafile.New2(c.clientProvider(), logger, nil)
}

// Config implements [di.Container].
func (c *Container) Config() domainconfig.ConfigService {
	c.configOnce.Do(func() {
		c.configService = c.newConfigService()
	})
	return c.configService
}

func (c *Container) newConfigService() domainconfig.ConfigService {
	loggerService := c.Logger()
	logger, err := loggerService.CreateLogger("config")
	if err != nil {
		logger = infralogging.NewNoopLogger()
	}

	return appconfig.New2(infraconfig.NewYAMLLoader(), logger)
}

func (c *Container) Cache() domaincache.Service2 {
	c.cacheOnce2.Do(func() {
		c.cacheService2 = c.newCacheService()
	})

	return c.cacheService2
}

func (c *Container) newCacheService() domaincache.Service2 {
	loggerService := c.Logger()
	logger, err := loggerService.CreateLogger("cache")
	if err != nil {
		logger = infralogging.NewNoopLogger()
	}

	return appcache.NewService2(logger)
}

// Auth implements [di.Container].
func (c *Container) Auth() domainauth.AuthService {
	c.authOnce.Do(func() {
		c.authService = c.newAuthService()
	})

	return c.authService
}

func (c *Container) newAuthService() domainauth.AuthService {
	credentialFactory := msal.NewMSALCredentialFactory()

	loggerService := c.Logger()
	logger, err := loggerService.CreateLogger("auth")
	if err != nil {
		logger = infralogging.NewNoopLogger()
	}

	return appauth.NewService2(c.authCache(), c.Config(), c.State(), logger, credentialFactory, c.Account())
}

// EnvironmentService implements [di.Container].
func (c *Container) EnvironmentService() domainenv.EnvironmentService {
	c.environmentOnce.Do(func() {
		c.environmentService = c.newEnvironmentService()
	})
	return c.environmentService
}

func (c *Container) newEnvironmentService() domainenv.EnvironmentService {
	svc := environment.New2("odc")
	_ = svc.EnsureAll()
	return svc
}

// FS implements [di.Container].
func (c *Container) FS() domainfs.Service {
	c.fsOnce.Do(func() {
		c.fsService = c.newFSService()
	})
	return c.fsService
}

func (c *Container) newFSService() domainfs.Service {
	loggerService := c.Logger()
	logger, _ := loggerService.CreateLogger("filesystem")

	resolver := appstate.NewDriverResolverAdapter(c.State())

	return appfs.NewService2(c.metadata(), c.contents(), resolver, logger)
}

func (c *Container) Account() domainaccount.Service {
	c.accountOnce.Do(func() {
		c.accountService = c.newAccountService()
	})
	return c.accountService
}

func (c *Container) newAccountService() domainaccount.Service {
	loggerSvc := c.Logger()
	logger, err := loggerSvc.CreateLogger("account")
	if err != nil {
		logger = infralogging.NewNoopLogger()
	}

	return appaccount.New(c.accountCache(), logger)
}

func (c *Container) Editor() domaineditor.Service {
	c.editorOnce.Do(func() {
		c.editorService = c.newEditorService()
	})
	return c.editorService
}

func (c *Container) newEditorService() domaineditor.Service {
	loggerSvc := c.Logger()
	logger, err := loggerSvc.CreateLogger("editor")
	if err != nil {
		logger = infralogging.NewNoopLogger()
	}

	launcher := infraeditor.NewLauncher(c.EnvironmentService(), logger)
	return appeditor.NewService(launcher, logger)
}

// Logger implements [di.Container].
func (c *Container) Logger() domainlogger.LoggerService {
	c.loggerOnce.Do(func() {
		c.loggerService = c.newLoggerService()
	})
	return c.loggerService
}

func (c *Container) newLoggerService() domainlogger.LoggerService {
	level, _ := c.EnvironmentService().LogLevel()

	opts := []domainlogger.Option{domainlogger.WithLogLevel(level), domainlogger.WithType(infralogging.TypeZap)}

	outputDest, _ := c.EnvironmentService().OutputDestination()
	switch outputDest {
	case infralogging.OutputDestinationFile:
		logHome, _ := c.EnvironmentService().LogDir()
		opts = append(opts, domainlogger.WithOutputDestinationFile(logHome))
	case infralogging.OutputDestinationStandardOut:
		opts = append(opts, domainlogger.WithOutputDestinationStandardOut())
	case infralogging.OutputDestinationStandardError:
		opts = append(opts, domainlogger.WithOutputDestinationStandardError())
	default:
	}

	svc, _ := applogging.NewLoggerService(opts...)
	svc.RegisterProvider(infralogging.TypeZap, infralogging.NewZapLoggerProvider())
	return svc
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

	loggerService := c.Logger()
	logger, err := loggerService.CreateLogger("profile")
	if err != nil {
		logger = infralogging.NewNoopLogger()
	}

	// Infra repository
	repo := infraprofile.NewFSProfileService(profileBaseDir)

	// App service (repo only)
	return appprofile.New(
		logger,
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
