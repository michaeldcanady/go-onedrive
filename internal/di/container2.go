// internal/di/container.go
package di

import (
	"context"
	"path/filepath"
	"sync"

	cacheservice "github.com/michaeldcanady/go-onedrive/internal/app/cache_service"
	cliprofileservicego "github.com/michaeldcanady/go-onedrive/internal/app/cli_profile_service.go"
	clientservice "github.com/michaeldcanady/go-onedrive/internal/app/client_service"
	configurationservice "github.com/michaeldcanady/go-onedrive/internal/app/configuration_service"
	credentialservice "github.com/michaeldcanady/go-onedrive/internal/app/credential_service"
	driveservice "github.com/michaeldcanady/go-onedrive/internal/app/drive_service"
	environmentservice "github.com/michaeldcanady/go-onedrive/internal/app/environment_service"
	loggerservice "github.com/michaeldcanady/go-onedrive/internal/app/logger_service"
	"github.com/michaeldcanady/go-onedrive/internal/event"

	"github.com/spf13/viper"
)

type Container1 struct {
	Ctx     context.Context
	Options RuntimeOptions

	// static
	EnvironmentService EnvironmentService
	LoggerService      LoggerService
	EventBus           EventBus

	// lazy
	cacheOnce    sync.Once
	cacheService CacheService
	cacheErr     error

	configOnce    sync.Once
	configService ConfigurationService
	configErr     error

	credOnce    sync.Once
	credService CredentialService
	credErr     error

	graphOnce    sync.Once
	graphService Clienter
	graphErr     error

	driveOnce    sync.Once
	driveService ChildrenIterator
	driveErr     error

	profileOnce    sync.Once
	profileService CLIProfileService
	profileErr     error
}

func NewContainer1(ctx context.Context) (*Container1, error) {
	c := &Container1{Ctx: ctx}

	// environment
	c.EnvironmentService = environmentservice.New("odc")

	// base logger (level may be overridden later via flags)
	logDir, err := c.EnvironmentService.LogDir(ctx)
	if err != nil {
		return nil, err
	}

	// default level; can be refined using viper in lazy paths
	c.LoggerService, err = loggerservice.New("info", logDir)
	if err != nil {
		return nil, err
	}

	// event bus
	busLogger, _ := c.LoggerService.CreateLogger("bus")
	c.EventBus = event.NewInMemoryBus(busLogger)

	return c, nil
}

func (c *Container1) ProfileService() (CLIProfileService, error) {
	c.profileOnce.Do(func() {
		logger, _ := c.LoggerService.CreateLogger("profile")

		cache, err := c.CacheService()
		if err != nil {
			c.configErr = err
			return
		}
		configDir, err := c.EnvironmentService.ConfigDir(c.Ctx)
		if err != nil {
			c.cacheErr = err
			return
		}

		c.profileService = cliprofileservicego.New(cache, logger, configDir)
	})

	return c.profileService, c.profileErr
}

// CacheService is lazy and flag‑aware if needed.
func (c *Container1) CacheService() (CacheService, error) {
	c.cacheOnce.Do(func() {
		logger, _ := c.LoggerService.CreateLogger("cache")

		cacheDir, err := c.EnvironmentService.CacheDir(c.Ctx)
		if err != nil {
			c.cacheErr = err
			return
		}

		path := filepath.Join(cacheDir, "profile.cache")
		c.cacheService, c.cacheErr = cacheservice.New(path, logger)
	})
	return c.cacheService, c.cacheErr
}

func (c *Container1) ConfigurationService() (ConfigurationService, error) {
	c.configOnce.Do(func() {
		cache, err := c.CacheService()
		if err != nil {
			c.configErr = err
			return
		}

		logger, _ := c.LoggerService.CreateLogger("configuration")

		svc := configurationservice.New2(cache, YAMLLoader{}, logger)

		// config path can come from env or flags (e.g. --config)
		// assume viper has been bound in root command
		configDir := c.Options.ConfigPath
		profileName := defaultProfileName
		if configDir == "" {
			profileService, err := c.ProfileService()
			if err != nil {
				c.configErr = err
				return
			}
			profile, err := profileService.GetProfile(context.Background(), c.Options.ProfileName)
			if err != nil {
				c.configErr = err
				return
			}
			profileName = c.Options.ProfileName
			configDir = profile.ConfigurationPath
		}

		svc.AddPath(profileName, configDir)
		c.configService = svc
	})
	return c.configService, c.configErr
}

func (c *Container1) CredentialService() (CredentialService, error) {
	c.credOnce.Do(func() {
		cache, err := c.CacheService()
		if err != nil {
			c.credErr = err
			return
		}

		configService, err := c.ConfigurationService()
		if err != nil {
			c.credErr = err
			return
		}

		logger, _ := c.LoggerService.CreateLogger("credential")
		c.credService = credentialservice.New(cache, c.EventBus, logger, configService)
	})
	return c.credService, c.credErr
}

func (c *Container1) GraphClientService() (Clienter, error) {
	c.graphOnce.Do(func() {
		cred, err := c.CredentialService()
		if err != nil {
			c.graphErr = err
			return
		}

		logger, _ := c.LoggerService.CreateLogger("graph")
		c.graphService = clientservice.New(cred, c.EventBus, logger)
	})
	return c.graphService, c.graphErr
}

func (c *Container1) DriveService() (ChildrenIterator, error) {
	c.driveOnce.Do(func() {
		graph, err := c.GraphClientService()
		if err != nil {
			c.driveErr = err
			return
		}

		logger, _ := c.LoggerService.CreateLogger("drive")

		// example: use --debug or logging.level
		if viper.GetBool("debug") {
			logger.SetLevel("debug")
		} else if lvl := viper.GetString("logging.level"); lvl != "" {
			logger.SetLevel(lvl)
		}

		c.driveService = driveservice.New(graph, c.EventBus, logger)
	})
	return c.driveService, c.driveErr
}
