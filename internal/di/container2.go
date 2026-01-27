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
	driveservice2 "github.com/michaeldcanady/go-onedrive/internal/app/drive_service2"
	environmentservice "github.com/michaeldcanady/go-onedrive/internal/app/environment_service"
	driveservice "github.com/michaeldcanady/go-onedrive/internal/app/file_service"
	loggerservice "github.com/michaeldcanady/go-onedrive/internal/app/logger_service"
	"github.com/michaeldcanady/go-onedrive/internal/event"

	"github.com/spf13/viper"
)

const (
	defaultProfileName = "default"
)

type Container1 struct {
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

	//DEPRECATED
	driveOnce sync.Once
	//DEPRECATED
	driveService ChildrenIterator
	//DEPRECATED
	driveErr error

	driveService2 DriveService
	driveOnce2    sync.Once
	driveErr2     error

	fileOnce    sync.Once
	fileService FileSystemService
	fileErr     error

	profileOnce    sync.Once
	profileService CLIProfileService
	profileErr     error
}

func NewContainer1() (*Container1, error) {
	c := &Container1{}

	// environment
	c.EnvironmentService = environmentservice.New("odc")

	// base logger (level may be overridden later via flags)
	logDir, err := c.EnvironmentService.LogDir(context.Background())
	if err != nil {
		return nil, err
	}

	// default level; can be refined using viper in lazy paths
	c.LoggerService, err = loggerservice.New("info", logDir, &loggerProvider{})
	if err != nil {
		return nil, err
	}

	// event bus
	busLogger, _ := c.LoggerService.CreateLogger("bus")
	c.EventBus = event.NewInMemoryBus(busLogger)

	return c, nil
}

func (c *Container1) ProfileService(ctx context.Context) (CLIProfileService, error) {
	c.profileOnce.Do(func() {
		logger, _ := c.LoggerService.CreateLogger("profile")

		cache, err := c.CacheService(ctx)
		if err != nil {
			c.configErr = err
			return
		}
		configDir, err := c.EnvironmentService.ConfigDir(ctx)
		if err != nil {
			c.cacheErr = err
			return
		}

		c.profileService = cliprofileservicego.New(cache, logger, configDir)
	})

	return c.profileService, c.profileErr
}

// CacheService is lazy and flag‑aware if needed.
func (c *Container1) CacheService(ctx context.Context) (CacheService, error) {
	c.cacheOnce.Do(func() {
		logger, _ := c.LoggerService.CreateLogger("cache")

		cacheDir, err := c.EnvironmentService.CacheDir(ctx)
		if err != nil {
			c.cacheErr = err
			return
		}

		profilePath := filepath.Join(cacheDir, "profile.cache")
		drivePath := filepath.Join(cacheDir, "drive.cache")
		filePath := filepath.Join(cacheDir, "file.cache")
		c.cacheService, c.cacheErr = cacheservice.New(profilePath, drivePath, filePath, logger)
	})
	return c.cacheService, c.cacheErr
}

func (c *Container1) ConfigurationService(ctx context.Context) (ConfigurationService, error) {
	c.configOnce.Do(func() {
		cache, err := c.CacheService(ctx)
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
			profileService, err := c.ProfileService(ctx)
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

func (c *Container1) CredentialService(ctx context.Context) (CredentialService, error) {
	c.credOnce.Do(func() {
		cache, err := c.CacheService(ctx)
		if err != nil {
			c.credErr = err
			return
		}

		configService, err := c.ConfigurationService(ctx)
		if err != nil {
			c.credErr = err
			return
		}

		logger, _ := c.LoggerService.CreateLogger("credential")
		c.credService = credentialservice.New(cache, c.EventBus, logger, configService)
	})
	return c.credService, c.credErr
}

func (c *Container1) GraphClientService(ctx context.Context) (Clienter, error) {
	c.graphOnce.Do(func() {
		cred, err := c.CredentialService(ctx)
		if err != nil {
			c.graphErr = err
			return
		}

		logger, _ := c.LoggerService.CreateLogger("graph")
		c.graphService = clientservice.New(cred, c.EventBus, logger)
	})
	return c.graphService, c.graphErr
}

func (c *Container1) FileService(ctx context.Context) (FileSystemService, error) {
	c.fileOnce.Do(func() {
		cache, err := c.CacheService(ctx)
		if err != nil {
			c.driveErr = err
			return
		}

		graph, err := c.GraphClientService(ctx)
		if err != nil {
			c.driveErr = err
			return
		}

		logger, _ := c.LoggerService.CreateLogger("file")

		c.fileService = driveservice.New2(graph, c.EventBus, logger, cache)
	})
	return c.fileService, c.fileErr
}

func (c *Container1) DriveService2(ctx context.Context) (DriveService, error) {
	c.driveOnce2.Do(func() {
		graph, err := c.GraphClientService(ctx)
		if err != nil {
			c.driveErr = err
			return
		}

		logger, _ := c.LoggerService.CreateLogger("drive2")

		// example: use --debug or logging.level
		if viper.GetBool("debug") {
			logger.SetLevel("debug")
		} else if lvl := viper.GetString("logging.level"); lvl != "" {
			logger.SetLevel(lvl)
		}

		c.driveService2 = driveservice2.NewDriveService(graph, logger)
	})
	return c.driveService2, c.driveErr2
}
