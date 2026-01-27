package di

import (
	"context"
	"path/filepath"
	"sync"

	cacheservice "github.com/michaeldcanady/go-onedrive/internal/app/cache_service"
	cliprofileservicego "github.com/michaeldcanady/go-onedrive/internal/app/cli_profile_service.go"
	configurationservice "github.com/michaeldcanady/go-onedrive/internal/app/configuration_service"
	driveservice2 "github.com/michaeldcanady/go-onedrive/internal/app/drive_service2"
	environmentservice "github.com/michaeldcanady/go-onedrive/internal/app/environment_service"
	loggerservice "github.com/michaeldcanady/go-onedrive/internal/app/logger_service"

	fileservice "github.com/michaeldcanady/go-onedrive/internal/app/file_service"
	"github.com/michaeldcanady/go-onedrive/internal/auth"
	"github.com/michaeldcanady/go-onedrive/internal/graph"

	"github.com/michaeldcanady/go-onedrive/internal/event"
	"github.com/spf13/viper"
)

type Container struct {
	Options RuntimeOptions

	// static
	EnvironmentService EnvironmentService
	LoggerService      LoggerService
	EventBus           EventBus

	// lazy
	cacheOnce sync.Once
	cache     CacheService
	cacheErr  error

	configOnce sync.Once
	configSvc  ConfigurationService
	configErr  error

	credOnce sync.Once
	credProv CredentialProvider
	credErr  error

	graphOnce sync.Once
	graphProv GraphClientProvider
	graphErr  error

	driveOnce sync.Once
	driveSvc  *driveservice2.Service
	driveErr  error

	fileOnce sync.Once
	fileSvc  *fileservice.Service2
	fileErr  error

	profileOnce sync.Once
	profileSvc  CLIProfileService
	profileErr  error
}

func NewContainer() (*Container, error) {
	c := &Container{}

	// environment
	c.EnvironmentService = environmentservice.New("odc")

	// base logger
	logDir, err := c.EnvironmentService.LogDir(context.Background())
	if err != nil {
		return nil, err
	}

	c.LoggerService, err = loggerservice.New("info", logDir, newZapLogger)
	if err != nil {
		return nil, err
	}

	// event bus
	busLogger, _ := c.LoggerService.CreateLogger("bus")
	c.EventBus = event.NewInMemoryBus(busLogger)

	return c, nil
}

func (c *Container) CacheService(ctx context.Context) (CacheService, error) {
	c.cacheOnce.Do(func() {
		logger, _ := c.LoggerService.GetLogger("cache")

		cacheDir, err := c.EnvironmentService.CacheDir(ctx)
		if err != nil {
			c.cacheErr = err
			return
		}

		profilePath := filepath.Join(cacheDir, "profile.cache")
		drivePath := filepath.Join(cacheDir, "drive.cache")
		filePath := filepath.Join(cacheDir, "file.cache")

		c.cache, c.cacheErr = cacheservice.New(profilePath, drivePath, filePath, logger)
	})
	return c.cache, c.cacheErr
}

func (c *Container) ConfigurationService(ctx context.Context) (ConfigurationService, error) {
	c.configOnce.Do(func() {
		cache, err := c.CacheService(ctx)
		if err != nil {
			c.configErr = err
			return
		}

		logger, _ := c.LoggerService.CreateLogger("configuration")
		svc := configurationservice.New2(cache, YAMLLoader{}, logger)

		configDir := c.Options.ConfigPath
		profileName := defaultProfileName

		if configDir == "" {
			profileSvc, err := c.ProfileService(ctx)
			if err != nil {
				c.configErr = err
				return
			}

			profile, err := profileSvc.GetProfile(ctx, c.Options.ProfileName)
			if err != nil {
				c.configErr = err
				return
			}

			profileName = c.Options.ProfileName
			configDir = profile.ConfigurationPath
		}

		svc.AddPath(profileName, configDir)
		c.configSvc = svc
	})
	return c.configSvc, c.configErr
}

func (c *Container) CredentialProvider(ctx context.Context) (CredentialProvider, error) {
	c.credOnce.Do(func() {
		cache, err := c.CacheService(ctx)
		if err != nil {
			c.credErr = err
			return
		}

		configSvc, err := c.ConfigurationService(ctx)
		if err != nil {
			c.credErr = err
			return
		}

		logger, _ := c.LoggerService.CreateLogger("credential")

		factory := auth.NewDefaultCredentialFactory()
		c.credProv = auth.NewCredentialProvider(cache, configSvc, factory, logger)
	})
	return c.credProv, c.credErr
}

func (c *Container) GraphClientProvider(ctx context.Context) (GraphClientProvider, error) {
	c.graphOnce.Do(func() {
		credProv, err := c.CredentialProvider(ctx)
		if err != nil {
			c.graphErr = err
			return
		}

		logger, _ := c.LoggerService.CreateLogger("graph")
		c.graphProv = graph.NewGraphClientProvider(credProv, logger)
	})
	return c.graphProv, c.graphErr
}

func (c *Container) DriveService(ctx context.Context) (*driveservice2.Service, error) {
	c.driveOnce.Do(func() {
		graphProv, err := c.GraphClientProvider(ctx)
		if err != nil {
			c.driveErr = err
			return
		}

		logger, _ := c.LoggerService.CreateLogger("drive")

		if viper.GetBool("debug") {
			logger.SetLevel("debug")
		} else if lvl := viper.GetString("logging.level"); lvl != "" {
			logger.SetLevel(lvl)
		}

		c.driveSvc = driveservice2.NewDriveService(graphProv, logger)
	})
	return c.driveSvc, c.driveErr
}

func (c *Container) FileSystemService(ctx context.Context) (*fileservice.Service2, error) {
	c.fileOnce.Do(func() {
		cache, err := c.CacheService(ctx)
		if err != nil {
			c.fileErr = err
			return
		}

		graphProv, err := c.GraphClientProvider(ctx)
		if err != nil {
			c.fileErr = err
			return
		}

		logger, _ := c.LoggerService.CreateLogger("file")
		c.fileSvc = fileservice.New2(graphProv, c.EventBus, logger, cache)
	})
	return c.fileSvc, c.fileErr
}

func (c *Container) ProfileService(ctx context.Context) (CLIProfileService, error) {
	c.profileOnce.Do(func() {
		logger, _ := c.LoggerService.CreateLogger("profile")

		cache, err := c.CacheService(ctx)
		if err != nil {
			c.profileErr = err
			return
		}

		configDir, err := c.EnvironmentService.ConfigDir(ctx)
		if err != nil {
			c.profileErr = err
			return
		}

		c.profileSvc = cliprofileservicego.New(cache, logger, configDir)
	})
	return c.profileSvc, c.profileErr
}
