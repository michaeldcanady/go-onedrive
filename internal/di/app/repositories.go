package app

import (
	"context"

	domaindrive "github.com/michaeldcanady/go-onedrive/internal/drive/domain"
	infradrive "github.com/michaeldcanady/go-onedrive/internal/drive/infra"
	apponedrive "github.com/michaeldcanady/go-onedrive/internal/fs/onedrive/app"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	infrafile "github.com/michaeldcanady/go-onedrive/internal/fs/onedrive/infra"
)

func (c *Container) metadata() domainfs.MetadataRepository {
	c.metadataOnce.Do(func() {
		c.metadataRepo = c.newMetadataRepository()
	})
	return c.metadataRepo
}

func (c *Container) newDriveRepository() domaindrive.DriveGateway {
	adapter, _ := c.clientProvider().RequestAdapter(context.Background())
	gateway := infradrive.NewGraphDriveGateway(adapter, c.getLogger("drive_gateway"))

	return gateway
}

func (c *Container) newMetadataRepository() domainfs.MetadataRepository {
	adapter, _ := c.clientProvider().RequestAdapter(context.Background())
	gateway := infrafile.NewGraphMetadataGateway(adapter, c.getLogger("gateway"))
	return apponedrive.NewCachedMetadataRepository(gateway, c.metadataCache(), c.metadataListingCache(), c.pathIDCache(), c.getLogger("repository"))
}

func (c *Container) contents() domainfs.FileContentsRepository {
	c.contentsOnce.Do(func() {
		c.contentsRepo = c.newContentsRepository()
	})
	return c.contentsRepo
}

func (c *Container) newContentsRepository() domainfs.FileContentsRepository {
	adapter, _ := c.clientProvider().RequestAdapter(context.Background())
	gateway := infrafile.NewGraphFileContentsGateway(adapter, c.getLogger("gateway"))
	return apponedrive.NewCachedFileContentsRepository(gateway, c.contentsCache(), c.metadataCache(), c.pathIDCache(), c.getLogger("repository"))
}
