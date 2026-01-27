package di2

import (
	"path/filepath"
	"sync"

	cacheservice "github.com/michaeldcanady/go-onedrive/internal/app/cache_service"
	cliprofileservicego "github.com/michaeldcanady/go-onedrive/internal/app/cli_profile_service.go"
	configurationservice "github.com/michaeldcanady/go-onedrive/internal/app/configuration_service"
	driveservice2 "github.com/michaeldcanady/go-onedrive/internal/app/drive_service2"
	environmentservice "github.com/michaeldcanady/go-onedrive/internal/app/environment_service"
	fileservice "github.com/michaeldcanady/go-onedrive/internal/app/file_service"
	loggerservice "github.com/michaeldcanady/go-onedrive/internal/app/logger_service"
	"github.com/michaeldcanady/go-onedrive/internal/auth"
	"github.com/michaeldcanady/go-onedrive/internal/event"
	"github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/graph"
)

type Container struct {
	once sync.Once

	env     EnvironmentService
	logger  LoggerService
	bus     EventBus
	cache   CacheService
	config  ConfigurationService
	creds   CredentialProvider
	graph   GraphClientProvider
	drives  DriveService
	files   FileService
	profile CLIProfileService

	fs fs.Service
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) FS() fs.Service {
	c.once.Do(func() {
		// environment
		c.env = environmentservice.New2("odc")

		// logger
		logDir, _ := c.env.LogDir()
		c.logger, _ = loggerservice.New("info", logDir, &loggerProvider{})

		// event bus
		busLogger, _ := c.logger.CreateLogger("bus")
		c.bus = event.NewInMemoryBus(busLogger)

		cacheLogger, _ := c.logger.CreateLogger("cache")

		// cache
		cacheDir, _ := c.env.CacheDir()
		c.cache, _ = cacheservice.New(
			filepath.Join(cacheDir, "profile.cache"),
			filepath.Join(cacheDir, "drive.cache"),
			filepath.Join(cacheDir, "file.cache"),
			cacheLogger,
		)

		configLogger, _ := c.logger.CreateLogger("config")

		// config
		c.config = configurationservice.New2(c.cache, YAMLLoader{}, configLogger)

		credentialLogger, _ := c.logger.CreateLogger("credential")

		// credentials
		factory := auth.NewDefaultCredentialFactory()
		c.creds = auth.NewCredentialProvider(c.cache, c.config, factory, credentialLogger)

		graphLogger, _ := c.logger.CreateLogger("graph")

		// graph
		c.graph = graph.NewGraphClientProvider(c.creds, graphLogger)

		driveLogger, _ := c.logger.CreateLogger("drive")

		// drive service
		c.drives = driveservice2.NewDriveService(c.graph, driveLogger)

		fileLogger, _ := c.logger.CreateLogger("file")

		// file service
		c.files = fileservice.New2(c.graph, c.bus, fileLogger, c.cache)

		profileLogger, _ := c.logger.CreateLogger("profile")

		// profile service
		configDir, _ := c.env.ConfigDir()
		c.profile = cliprofileservicego.New(c.cache, profileLogger, configDir)

		// filesystem service
		dr := &fs.PersonalDriveResolver{DriveService: c.drives}
		c.fs = fs.NewService(c.files, dr)
	})

	return c.fs
}
