package app

import (
	"sync"

	domainaccount "github.com/michaeldcanady/go-onedrive/internal/account/domain"
	domainauth "github.com/michaeldcanady/go-onedrive/internal/auth/domain"
	domainconfig "github.com/michaeldcanady/go-onedrive/internal/config/domain"
	domainenv "github.com/michaeldcanady/go-onedrive/internal/core/env/domain"
	domaingraph "github.com/michaeldcanady/go-onedrive/internal/core/graph/domain"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
	domaindrive "github.com/michaeldcanady/go-onedrive/internal/drive/domain"
	domaineditor "github.com/michaeldcanady/go-onedrive/internal/editor/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/shared/domain"
	domainprofile "github.com/michaeldcanady/go-onedrive/internal/profile/domain"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
	pkgcache "github.com/michaeldcanady/go-onedrive/pkg/cache"
	infrabolt "github.com/michaeldcanady/go-onedrive/pkg/cache/bolt"

	infrafile "github.com/michaeldcanady/go-onedrive/internal/fs/onedrive/infra"
)

var _ didomain.Container = (*Container)(nil)

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
	sharedStore    *infrabolt.Store

	cacheOnce2    sync.Once
	cacheService2 pkgcache.Service

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
	metadataRepo domainfs.MetadataRepository

	contentsOnce sync.Once
	contentsRepo domainfs.FileContentsRepository

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
