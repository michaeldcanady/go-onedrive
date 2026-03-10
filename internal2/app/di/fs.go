package di

import (
	appcache "github.com/michaeldcanady/go-onedrive/internal2/app/cache"
	appdrive "github.com/michaeldcanady/go-onedrive/internal2/app/drive"
	appeditor "github.com/michaeldcanady/go-onedrive/internal2/app/editor"
	appfs "github.com/michaeldcanady/go-onedrive/internal2/app/fs"
	applocal "github.com/michaeldcanady/go-onedrive/internal2/app/fs/local"
	apponedrive "github.com/michaeldcanady/go-onedrive/internal2/app/fs/onedrive"
	appregistry "github.com/michaeldcanady/go-onedrive/internal2/app/fs/registry"
	appstate "github.com/michaeldcanady/go-onedrive/internal2/app/state"

	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	domaindrive "github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	domaineditor "github.com/michaeldcanady/go-onedrive/internal2/domain/editor"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"

	infraeditor "github.com/michaeldcanady/go-onedrive/internal2/infra/editor"
)

// Drive implements [di.Container].
func (c *Container) Drive() domaindrive.DriveService {
	c.driveOnce.Do(func() {
		c.driveService = c.newDriveService()
	})
	return c.driveService
}

func (c *Container) newDriveService() domaindrive.DriveService {
	return appdrive.NewDriveService(c.clientProvider(), c.getLogger("drive"))
}

func (c *Container) Cache() domaincache.Service2 {
	c.cacheOnce2.Do(func() {
		c.cacheService2 = c.newCacheService()
	})

	return c.cacheService2
}

func (c *Container) newCacheService() domaincache.Service2 {
	return appcache.NewService2(c.getLogger("cache"))
}

// FS implements [di.Container].
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
