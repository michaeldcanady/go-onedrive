package di

import (
	"context"
	"path"
	"path/filepath"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal2/app/account"
	"github.com/michaeldcanady/go-onedrive/internal2/app/auth"
	"github.com/michaeldcanady/go-onedrive/internal2/app/cache"
	"github.com/michaeldcanady/go-onedrive/internal2/app/config"
	"github.com/michaeldcanady/go-onedrive/internal2/app/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/app/fs"
	appprofile "github.com/michaeldcanady/go-onedrive/internal2/app/profile"
	"github.com/michaeldcanady/go-onedrive/internal2/app/state"

	domainaccount "github.com/michaeldcanady/go-onedrive/internal2/domain/account"

	"github.com/michaeldcanady/go-onedrive/internal2/app/common/environment"
	"github.com/michaeldcanady/go-onedrive/internal2/app/common/logging"
	domainauth "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	domainenv "github.com/michaeldcanady/go-onedrive/internal2/domain/common/environment"
	domaingraph "github.com/michaeldcanady/go-onedrive/internal2/domain/common/graph"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	domainconfig "github.com/michaeldcanady/go-onedrive/internal2/domain/config"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	domaindrive "github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	domainfile "github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	domainprofile "github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
	domainstate "github.com/michaeldcanady/go-onedrive/internal2/domain/state"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/auth/msal"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/graph"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/file"
	infraprofile "github.com/michaeldcanady/go-onedrive/internal2/infra/profile"

	appstate "github.com/michaeldcanady/go-onedrive/internal2/app/state"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	infraconfig "github.com/michaeldcanady/go-onedrive/internal2/infra/config"
	infrastate "github.com/michaeldcanady/go-onedrive/internal2/infra/state"
)

var _ di.Container = (*Container)(nil)

const (
	stateFileName = "state.json"
)

type Container struct {
	authOnce    sync.Once
	authService domainauth.AuthService

	environmentOnce    sync.Once
	environmentService domainenv.EnvironmentService

	fsOnce    sync.Once
	fsService domainfs.Service

	loggerOnce    sync.Once
	loggerService domainlogger.LoggerService

	profileOnce    sync.Once
	profileService domainprofile.ProfileService

	cacheOnce2    sync.Once
	cacheService2 domaincache.Service2

	// DEPRECATED: use cacheOnce2 instead.
	cacheOnce sync.Once

	// DEPRECATED: use cacheService2 instead.
	cacheService domaincache.CacheService

	configOnce    sync.Once
	configService domainconfig.ConfigService

	// DEPRECATED
	fileOnce sync.Once
	// DEPRECATED
	fileService domainfile.FileService

	clientOnce    sync.Once
	clientProvide domaingraph.ClientProvider

	stateOnce    sync.Once
	stateService domainstate.Service

	driveOnce    sync.Once
	driveService domaindrive.DriveService

	accountOnce    sync.Once
	accountService domainaccount.Service

	metadataOnce sync.Once
	metadataRepo *file.MetadataRepository

	contentsOnce sync.Once
	contentsRepo *file.ContentsRepository

	metadataCacheOnce  sync.Once
	metadataCacheCache file.MetadataCache

	metadataListingCacheOnce  sync.Once
	metadataListingCacheCache file.ListingCache

	contentsCacheOnce  sync.Once
	contentsCacheCache file.ContentsCache
}

func NewContainer() *Container {
	return &Container{}
}

// Drive implements [di.Container].
func (c *Container) Drive() domaindrive.DriveService {
	c.driveOnce.Do(func() {
		loggerService := c.Logger()
		logger, _ := loggerService.CreateLogger("drive")

		c.driveService = drive.NewDriveService(c.clientProvider(), logger)
	})

	return c.driveService
}

// File implements [di.Container].
func (c *Container) File() domainfile.FileService {
	c.fileOnce.Do(func() {
		loggerService := c.Logger()
		logger, _ := loggerService.CreateLogger("file")

		c.fileService = file.New2(c.clientProvider(), logger, c.CacheService())
	})
	return c.fileService
}

// Config implements [di.Container].
func (c *Container) Config() domainconfig.ConfigService {
	c.configOnce.Do(func() {
		loggerService := c.Logger()
		logger, _ := loggerService.CreateLogger("config")

		c.configService = config.New2(c.CacheService(), infraconfig.NewYAMLLoader(), logger)
	})
	return c.configService
}

func (c *Container) Cache() domaincache.Service2 {
	c.cacheOnce2.Do(func() {
		loggerService := c.Logger()
		logger, _ := loggerService.CreateLogger("cache")

		c.cacheService2 = cache.NewService2(logger)
	})

	return c.cacheService2
}

// DEPRECATED: use Cache instead.
func (c *Container) CacheService() domaincache.CacheService {
	c.cacheOnce.Do(func() {
		environmentService := c.EnvironmentService()
		cachePath, _ := environmentService.CacheDir()

		c.cacheService = cache.NewServiceAdapter(filepath.Join(cachePath, "cache.db"), c.Cache())
	})
	return c.cacheService
}

// Auth implements [di.Container].
func (c *Container) Auth() domainauth.AuthService {
	c.authOnce.Do(func() {
		credentialFactory := msal.NewMSALCredentialFactory()

		loggerService := c.Logger()
		logger, _ := loggerService.CreateLogger("auth")

		environmentService := c.EnvironmentService()
		cachePath, _ := environmentService.CacheDir()

		cacheSvc := c.Cache()
		cache := cacheSvc.CreateCache(context.Background(), "auth", cache.BoltCacheFactory(path.Join(cachePath, "auth.db"), "auth"))

		c.authService = auth.NewService2(cache, c.Config(), c.State(), logger, credentialFactory, c.Account())
	})

	return c.authService
}

func (c *Container) clientProvider() domaingraph.ClientProvider {
	c.clientOnce.Do(func() {
		loggerService := c.Logger()
		graphLogger, _ := loggerService.CreateLogger("graph")
		c.clientProvide = graph.New(c.Auth(), graphLogger)
	})

	return c.clientProvide
}

// EnvironmentService implements [di.Container].
func (c *Container) EnvironmentService() domainenv.EnvironmentService {
	c.environmentOnce.Do(func() {
		c.environmentService = environment.New2("odc")

		_ = c.environmentService.EnsureAll()
	})
	return c.environmentService
}

func (c *Container) metadataCache() file.MetadataCache {
	c.metadataCacheOnce.Do(func() {
		environmentService := c.EnvironmentService()
		cachePath, _ := environmentService.CacheDir()

		store := cache.BoltCacheFactory(path.Join(cachePath, "metadata.db"), "metadata")

		cache := c.Cache().CreateCache(context.Background(), "metadata", store)

		c.metadataCacheCache = file.NewMetadataCacheAdapter(cache)
	})
	return c.metadataCacheCache
}

func (c *Container) metadataListingCache() file.ListingCache {
	c.metadataListingCacheOnce.Do(func() {
		environmentService := c.EnvironmentService()
		cachePath, _ := environmentService.CacheDir()

		store := cache.BoltCacheFactory(path.Join(cachePath, "metadatl.db"), "metadatl")

		cache := c.Cache().CreateCache(context.Background(), "metadatl", store)

		c.metadataListingCacheCache = file.NewMetadataListCacheAdapter(cache)
	})
	return c.metadataListingCacheCache
}

func (c *Container) contentsCache() file.ContentsCache {
	c.metadataCacheOnce.Do(func() {
		environmentService := c.EnvironmentService()
		cachePath, _ := environmentService.CacheDir()

		store := cache.BoltCacheFactory(path.Join(cachePath, "contents.db"), "contents")

		cache := c.Cache().CreateCache(context.Background(), "contents", store)

		c.contentsCacheCache = file.NewContentsCacheAdapter(cache)
	})
	return c.contentsCacheCache
}

func (c *Container) metadata() *file.MetadataRepository {
	c.metadataOnce.Do(func() {

		client, _ := c.clientProvider().Client(context.Background())

		c.metadataRepo = file.NewMetadataRepository(client.RequestAdapter, c.metadataCache(), c.metadataListingCache())
	})
	return c.metadataRepo
}

func (c *Container) contents() *file.ContentsRepository {
	c.contentsOnce.Do(func() {
		client, _ := c.clientProvider().Client(context.Background())

		c.contentsRepo = file.NewContentsRepository(client.RequestAdapter, c.contentsCache(), c.metadataCache())
	})
	return c.contentsRepo
}

// FS implements [di.Container].
func (c *Container) FS() domainfs.Service {
	c.fsOnce.Do(func() {
		loggerService := c.Logger()
		logger, _ := loggerService.CreateLogger("filesystem")

		resolver := state.NewDriverResolverAdapter(c.State())

		c.fsService = fs.NewService2(c.metadata(), c.contents(), resolver, logger)
	})
	return c.fsService
}

func (c *Container) Account() domainaccount.Service {
	c.accountOnce.Do(func() {
		cacheSvc := c.Cache()

		loggerSvc := c.Logger()
		logger, _ := loggerSvc.CreateLogger("account")

		environmentService := c.EnvironmentService()
		cachePath, _ := environmentService.CacheDir()

		cache := cacheSvc.CreateCache(context.Background(), "account", cache.BoltCacheFactory(path.Join(cachePath, "account.db"), "account"))

		c.accountService = account.New(cache, logger)
	})
	return c.accountService
}

// Logger implements [di.Container].
func (c *Container) Logger() domainlogger.LoggerService {
	c.loggerOnce.Do(func() {
		level, _ := c.EnvironmentService().LogLevel()

		opts := []logger.Option{logger.WithLogLevel(level), logger.WithType(infralogging.TypeZap)}

		outputDest, _ := c.EnvironmentService().OutputDestination()
		switch outputDest {
		case infralogging.OutputDestinationFile:
			logHome, _ := c.EnvironmentService().LogDir()
			opts = append(opts, logger.WithOutputDestinationFile(logHome))
		case infralogging.OutputDestinationStandardOut:
			opts = append(opts, logger.WithOutputDestinationStandardOut())
		case infralogging.OutputDestinationStandardError:
			opts = append(opts, logger.WithOutputDestinationStandardError())
		default:
		}

		c.loggerService, _ = logging.NewLoggerService(opts...)
		c.loggerService.RegisterProvider(infralogging.TypeZap, infralogging.NewZapLoggerProvider())
	})
	return c.loggerService
}

// Profile implements [di.Container].
func (c *Container) Profile() domainprofile.ProfileService {
	c.profileOnce.Do(func() {
		env := c.EnvironmentService()

		// ~/.config/odc
		profileBaseDir, _ := env.ConfigDir()

		loggerService := c.Logger()
		logger, _ := loggerService.CreateLogger("profile")

		// Infra repository
		repo := infraprofile.NewFSProfileService(profileBaseDir)

		// App service (cache + repo)
		c.profileService = appprofile.New(
			c.CacheService(),
			logger,
			repo,
		)
	})

	return c.profileService
}

func (c *Container) State() domainstate.Service {
	c.stateOnce.Do(func() {
		env := c.EnvironmentService()
		stateDir, _ := env.StateDir()
		statePath := filepath.Join(stateDir, stateFileName)

		serializer := &cache.JSONSerializerDeserializer[domainstate.State]{}
		repo := infrastate.NewRepository(statePath, serializer)

		c.stateService = appstate.NewService(repo)
	})
	return c.stateService
}
