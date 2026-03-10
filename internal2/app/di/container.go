package di

import (
	"sync"

	domainaccount "github.com/michaeldcanady/go-onedrive/internal2/domain/account"
	domainauth "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	domainenv "github.com/michaeldcanady/go-onedrive/internal2/domain/common/environment"
	domaingraph "github.com/michaeldcanady/go-onedrive/internal2/domain/common/graph"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	domainconfig "github.com/michaeldcanady/go-onedrive/internal2/domain/config"
	domaindi "github.com/michaeldcanady/go-onedrive/internal2/domain/di"
	domaindrive "github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	domaineditor "github.com/michaeldcanady/go-onedrive/internal2/domain/editor"
	domainfile "github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	domainprofile "github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
	domainstate "github.com/michaeldcanady/go-onedrive/internal2/domain/state"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/bolt"
	infrafile "github.com/michaeldcanady/go-onedrive/internal2/infra/file"
)

var _ domaindi.Container = (*Container)(nil)

const (
	stateFileName = "state.json"
)

// Container implements the dependency injection container for the application.
// It uses a lazy initialization pattern with sync.Once.
type Container struct {
	// Core Services
	environmentOnce    sync.Once
	environmentService domainenv.EnvironmentService

	loggerOnce    sync.Once
	loggerService domainlogger.LoggerService

	stateOnce    sync.Once
	stateService domainstate.Service

	configOnce    sync.Once
	configService domainconfig.ConfigService

	// Infrastructure & Shared Components
	clientOnce    sync.Once
	clientProvide domaingraph.ClientProvider

	cacheStoreOnce sync.Once
	sharedStore    *bolt.Store

	cacheOnce2    sync.Once
	cacheService2 domaincache.Service2

	// Domain Services
	authOnce    sync.Once
	authService domainauth.AuthService

	accountOnce    sync.Once
	accountService domainaccount.Service

	profileOnce    sync.Once
	profileService domainprofile.ProfileService

	driveOnce    sync.Once
	driveService domaindrive.DriveService

	fsOnce    sync.Once
	fsService domainfs.Service

	ignoreMatcherFactoryOnce    sync.Once
	ignoreMatcherFactoryService domainfs.IgnoreMatcherFactory

	editorOnce    sync.Once
	editorService domaineditor.Service

	// Repository Components
	metadataOnce sync.Once
	metadataRepo domainfile.MetadataRepository

	contentsOnce sync.Once
	contentsRepo domainfile.FileContentsRepository

	// Caches
	metadataCacheOnce  sync.Once
	metadataCacheCache infrafile.MetadataCache

	metadataListingCacheOnce  sync.Once
	metadataListingCacheCache infrafile.ListingCache

	contentsCacheOnce  sync.Once
	contentsCacheCache infrafile.ContentsCache
}

// NewContainer creates a new instance of the dependency injection container.
func NewContainer() *Container {
	return &Container{}
}
