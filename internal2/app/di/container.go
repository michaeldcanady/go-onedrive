package di

import (
	"context"
	"path/filepath"
	"sync"

	appaccount "github.com/michaeldcanady/go-onedrive/internal2/app/account"
	appauth "github.com/michaeldcanady/go-onedrive/internal2/app/auth"
	appcache "github.com/michaeldcanady/go-onedrive/internal2/app/cache"
	appconfig "github.com/michaeldcanady/go-onedrive/internal2/app/config"
	appdrive "github.com/michaeldcanady/go-onedrive/internal2/app/drive"
	appfs "github.com/michaeldcanady/go-onedrive/internal2/app/fs"
	applogging "github.com/michaeldcanady/go-onedrive/internal2/app/common/logging"
	appprofile "github.com/michaeldcanady/go-onedrive/internal2/app/profile"
	appstate "github.com/michaeldcanady/go-onedrive/internal2/app/state"

	domainaccount "github.com/michaeldcanady/go-onedrive/internal2/domain/account"
	domainauth "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	domainenv "github.com/michaeldcanady/go-onedrive/internal2/domain/common/environment"
	domaingraph "github.com/michaeldcanady/go-onedrive/internal2/domain/common/graph"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	domainconfig "github.com/michaeldcanady/go-onedrive/internal2/domain/config"
	domaindi "github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	domaindrive "github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	domainfile "github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	domainprofile "github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
	domainstate "github.com/michaeldcanady/go-onedrive/internal2/domain/state"

	"github.com/michaeldcanady/go-onedrive/internal2/app/common/environment"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/auth/msal"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/bolt"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/graph"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	infraconfig "github.com/michaeldcanady/go-onedrive/internal2/infra/config"
	infrafile "github.com/michaeldcanady/go-onedrive/internal2/infra/file"
	infraprofile "github.com/michaeldcanady/go-onedrive/internal2/infra/profile"
	infrastate "github.com/michaeldcanady/go-onedrive/internal2/infra/state"
)

var _ domaindi.Container = (*Container)(nil)

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
	metadataRepo *infrafile.MetadataRepository

	contentsOnce sync.Once
	contentsRepo *infrafile.ContentsRepository

	metadataCacheOnce  sync.Once
	metadataCacheCache infrafile.MetadataCache

	metadataListingCacheOnce  sync.Once
	metadataListingCacheCache infrafile.ListingCache

	contentsCacheOnce  sync.Once
	contentsCacheCache infrafile.ContentsCache

	cacheStoreOnce sync.Once
	sharedStore    *bolt.Store
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) cacheStore() *bolt.Store {
	c.cacheStoreOnce.Do(func() {
		environmentService := c.EnvironmentService()
		cachePath, _ := environmentService.CacheDir()
		dbPath := filepath.Join(cachePath, "cache.db")

		// Create a "root" store with a dummy bucket just to open the DB
		store, err := bolt.NewStore(dbPath, "root")
		if err != nil {
			panic(err) // Critical failure if we can't open the cache
		}
		c.sharedStore = store
	})
	return c.sharedStore
}

// Drive implements [di.Container].
func (c *Container) Drive() domaindrive.DriveService {
	c.driveOnce.Do(func() {
		loggerService := c.Logger()
		logger, err := loggerService.CreateLogger("drive")
		if err != nil {
			// Fallback: use a default logger if creation fails, or at least note it.
			// Since we can't return an error from Once.Do easily without complex state, 
			// we'll use the service-level logger as fallback.
			logger = infralogging.NewNoopLogger()
		}

		c.driveService = appdrive.NewDriveService(c.clientProvider(), logger)
	})

	return c.driveService
}

// File implements [di.Container].
func (c *Container) File() domainfile.FileService {
	c.fileOnce.Do(func() {
		loggerService := c.Logger()
		logger, err := loggerService.CreateLogger("file")
		if err != nil {
			logger = infralogging.NewNoopLogger()
		}

		c.fileService = infrafile.New2(c.clientProvider(), logger, nil)
	})
	return c.fileService
}

// Config implements [di.Container].
func (c *Container) Config() domainconfig.ConfigService {
	c.configOnce.Do(func() {
		loggerService := c.Logger()
		logger, err := loggerService.CreateLogger("config")
		if err != nil {
			logger = infralogging.NewNoopLogger()
		}

		c.configService = appconfig.New2(infraconfig.NewYAMLLoader(), logger)
	})
	return c.configService
}

func (c *Container) Cache() domaincache.Service2 {
	c.cacheOnce2.Do(func() {
		loggerService := c.Logger()
		logger, err := loggerService.CreateLogger("cache")
		if err != nil {
			logger = infralogging.NewNoopLogger()
		}

		c.cacheService2 = appcache.NewService2(logger)
	})

	return c.cacheService2
}

func (c *Container) authCache() domaincache.Cache[domainauth.AccessToken] {
	cacheSvc := c.Cache()
	rawCache := cacheSvc.CreateCache(context.Background(), "auth_tokens", appcache.SiblingBoltFactory(c.cacheStore(), "auth_tokens"))

	return appcache.NewTypedCache(rawCache, &appcache.JSONSerializerDeserializer[domainauth.AccessToken]{})
}

// Auth implements [di.Container].
func (c *Container) Auth() domainauth.AuthService {
	c.authOnce.Do(func() {
		credentialFactory := msal.NewMSALCredentialFactory()

		loggerService := c.Logger()
		logger, err := loggerService.CreateLogger("auth")
		if err != nil {
			logger = infralogging.NewNoopLogger()
		}

		c.authService = appauth.NewService2(c.authCache(), c.Config(), c.State(), logger, credentialFactory, c.Account())
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

func (c *Container) metadataCache() infrafile.MetadataCache {
	c.metadataCacheOnce.Do(func() {
		cacheSvc := c.Cache()
		rawCache := cacheSvc.CreateCache(context.Background(), "metadata", appcache.SiblingBoltFactory(c.cacheStore(), "metadata"))

		typedCache := appcache.NewTypedCache(rawCache, &appcache.JSONSerializerDeserializer[domainfile.Metadata]{})
		c.metadataCacheCache = infrafile.NewMetadataCacheAdapter(typedCache)
	})
	return c.metadataCacheCache
}

func (c *Container) metadataListingCache() infrafile.ListingCache {
	c.metadataListingCacheOnce.Do(func() {
		cacheSvc := c.Cache()
		rawCache := cacheSvc.CreateCache(context.Background(), "metadatl", appcache.SiblingBoltFactory(c.cacheStore(), "metadatl"))

		typedCache := appcache.NewTypedCache(rawCache, &appcache.JSONSerializerDeserializer[infrafile.Listing]{})
		c.metadataListingCacheCache = infrafile.NewMetadataListCacheAdapter(typedCache)
	})
	return c.metadataListingCacheCache
}

func (c *Container) contentsCache() infrafile.ContentsCache {
	c.contentsCacheOnce.Do(func() {
		cacheSvc := c.Cache()
		rawCache := cacheSvc.CreateCache(context.Background(), "contents", appcache.SiblingBoltFactory(c.cacheStore(), "contents"))

		typedCache := appcache.NewTypedCache(rawCache, &appcache.JSONSerializerDeserializer[domainfile.Contents]{})
		c.contentsCacheCache = infrafile.NewContentsCacheAdapter(typedCache)
	})
	return c.contentsCacheCache
}

func (c *Container) pathIDCache() infrafile.PathIDCache {
	cacheSvc := c.Cache()
	rawCache := cacheSvc.CreateCache(context.Background(), "path_id", appcache.SiblingBoltFactory(c.cacheStore(), "path_id"))

	typedCache := appcache.NewTypedCache(rawCache, &appcache.JSONSerializerDeserializer[string]{})
	return infrafile.NewPathIDCacheAdapter(typedCache)
}

func (c *Container) metadata() *infrafile.MetadataRepository {
	c.metadataOnce.Do(func() {
		loggerSvc := c.Logger()
		logger, err := loggerSvc.CreateLogger("repository")
		if err != nil {
			logger = infralogging.NewNoopLogger()
		}

		client, _ := c.clientProvider().Client(context.Background())

		c.metadataRepo = infrafile.NewMetadataRepository(client.RequestAdapter, c.metadataCache(), c.metadataListingCache(), c.pathIDCache(), logger)
	})
	return c.metadataRepo
}

func (c *Container) contents() *infrafile.ContentsRepository {
	c.contentsOnce.Do(func() {
		loggerSvc := c.Logger()
		logger, err := loggerSvc.CreateLogger("repository")
		if err != nil {
			logger = infralogging.NewNoopLogger()
		}

		client, _ := c.clientProvider().Client(context.Background())

		c.contentsRepo = infrafile.NewContentsRepository(client.RequestAdapter, c.contentsCache(), c.metadataCache(), c.pathIDCache(), logger)
	})
	return c.contentsRepo
}

// FS implements [di.Container].
func (c *Container) FS() domainfs.Service {
	c.fsOnce.Do(func() {
		loggerService := c.Logger()
		logger, _ := loggerService.CreateLogger("filesystem")

		resolver := appstate.NewDriverResolverAdapter(c.State())

		c.fsService = appfs.NewService2(c.metadata(), c.contents(), resolver, logger)
	})
	return c.fsService
}

func (c *Container) accountCache() domaincache.Cache[domainaccount.Account] {
	cacheSvc := c.Cache()
	rawCache := cacheSvc.CreateCache(context.Background(), "account", appcache.SiblingBoltFactory(c.cacheStore(), "account"))

	return appcache.NewTypedCache(rawCache, &appcache.JSONSerializerDeserializer[domainaccount.Account]{})
}

func (c *Container) Account() domainaccount.Service {
	c.accountOnce.Do(func() {
		loggerSvc := c.Logger()
		logger, err := loggerSvc.CreateLogger("account")
		if err != nil {
			logger = infralogging.NewNoopLogger()
		}

		c.accountService = appaccount.New(c.accountCache(), logger)
	})
	return c.accountService
}

// Logger implements [di.Container].
func (c *Container) Logger() domainlogger.LoggerService {
	c.loggerOnce.Do(func() {
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

		c.loggerService, _ = applogging.NewLoggerService(opts...)
		c.loggerService.RegisterProvider(infralogging.TypeZap, infralogging.NewZapLoggerProvider())
	})
	return c.loggerService
}

// Profile implements [di.Container].
func (c *Container) Profile() domainprofile.ProfileService {
	c.profileOnce.Do(func() {
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
		c.profileService = appprofile.New(
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

		serializer := &appcache.JSONSerializerDeserializer[domainstate.State]{}
		repo := infrastate.NewRepository(statePath, serializer)

		c.stateService = appstate.NewService(repo)
	})
	return c.stateService
}
