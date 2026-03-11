package app

import (
	appcache "github.com/michaeldcanady/go-onedrive/internal/cache/app"
	appdrive "github.com/michaeldcanady/go-onedrive/internal/drive/app"
	appeditor "github.com/michaeldcanady/go-onedrive/internal/editor/app"
	applocal "github.com/michaeldcanady/go-onedrive/internal/fs/local/app"
	apponedrive "github.com/michaeldcanady/go-onedrive/internal/fs/onedrive/app"
	appregistry "github.com/michaeldcanady/go-onedrive/internal/fs/registry"
	appfs "github.com/michaeldcanady/go-onedrive/internal/fs/shared/app"
	appstate "github.com/michaeldcanady/go-onedrive/internal/state/app"

	domaindrive "github.com/michaeldcanady/go-onedrive/internal/drive/domain"
	domaineditor "github.com/michaeldcanady/go-onedrive/internal/editor/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/shared/domain"
	pkgcache "github.com/michaeldcanady/go-onedrive/pkg/cache"

	infraeditor "github.com/michaeldcanady/go-onedrive/internal/editor/infra"
)

// Drive implements [didomain.Container].
func (c *Container) Drive() domaindrive.DriveService {
	c.driveOnce.Do(func() {
		c.driveService = c.newDriveService()
	})
	return c.driveService
}

func (c *Container) newDriveService() domaindrive.DriveService {
	return appdrive.NewDriveService(c.newDriveRepository(), c.getLogger("drive"))
}

func (c *Container) Cache() pkgcache.Service {
	c.cacheOnce2.Do(func() {
		c.cacheService2 = c.newCacheService()
	})

	return c.cacheService2
}

func (c *Container) newCacheService() pkgcache.Service {
	return appcache.NewService2(c.getLogger("cache"))
}

// FS implements [didomain.Container].
func (c *Container) FS() domainfs.Service {
	c.fsOnce.Do(func() {
		c.fsService = c.newFSService()
	})
	return c.fsService
}

func (c *Container) newFSService() domainfs.Service {
	resolver := appstate.NewDriverResolverAdapter(c.State(), c.Drive())
	aliasSvc := appdrive.NewAliasService(c.State())

	// Create registry
	reg := appregistry.NewRegistry()

	// Register OneDrive provider
	oneDriveProvider := apponedrive.NewProvider(c.metadata(), c.contents(), resolver, aliasSvc, c.getLogger("filesystem"))
	reg.Register("onedrive", oneDriveProvider)

	// Register Local provider
	localProvider := applocal.NewProvider(c.getLogger("localfs"))
	reg.Register("local", localProvider)

	// Create manager
	return appfs.NewFileSystemManager(reg)
}

func (c *Container) IgnoreMatcherFactory() domainfs.IgnoreMatcherFactory {
	c.ignoreMatcherFactoryOnce.Do(func() {
		c.ignoreMatcherFactoryService = c.newIgnoreMatcherFactory()
	})
	return c.ignoreMatcherFactoryService
}

func (c *Container) newIgnoreMatcherFactory() domainfs.IgnoreMatcherFactory {
	return apponedrive.NewIgnoreMatcherFactory()
}

func (c *Container) Editor() domaineditor.Service {
	c.editorOnce.Do(func() {
		c.editorService = c.newEditorService()
	})
	return c.editorService
}

func (c *Container) newEditorService() domaineditor.Service {
	log := c.getLogger("editor")
	launcher := infraeditor.NewLauncher(c.EnvironmentService(), log)
	return appeditor.NewService(launcher, log)
}
