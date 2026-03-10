package di

import (
	"context"

	apponedrive "github.com/michaeldcanady/go-onedrive/internal2/app/fs/onedrive"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	infradrive "github.com/michaeldcanady/go-onedrive/internal2/infra/drive"
	infrafile "github.com/michaeldcanady/go-onedrive/internal2/infra/file"
)

func (c *Container) metadata() file.MetadataRepository {
	c.metadataOnce.Do(func() {
		c.metadataRepo = c.newMetadataRepository()
	})
	return c.metadataRepo
}

func (c *Container) newDriveRepository() drive.DriveGateway {
	adapter, _ := c.clientProvider().RequestAdapter(context.Background())
	gateway := infradrive.NewGraphDriveGateway(adapter, c.getLogger("drive_gateway"))

	return gateway
}

func (c *Container) newMetadataRepository() file.MetadataRepository {
	adapter, _ := c.clientProvider().RequestAdapter(context.Background())
	gateway := infrafile.NewGraphMetadataGateway(adapter, c.getLogger("gateway"))
	return apponedrive.NewCachedMetadataRepository(gateway, c.metadataCache(), c.metadataListingCache(), c.pathIDCache(), c.getLogger("repository"))
}

func (c *Container) contents() file.FileContentsRepository {
	c.contentsOnce.Do(func() {
		c.contentsRepo = c.newContentsRepository()
	})
	return c.contentsRepo
}

func (c *Container) newContentsRepository() file.FileContentsRepository {
	adapter, _ := c.clientProvider().RequestAdapter(context.Background())
	gateway := infrafile.NewGraphFileContentsGateway(adapter, c.getLogger("gateway"))
	return apponedrive.NewCachedFileContentsRepository(gateway, c.contentsCache(), c.metadataCache(), c.pathIDCache(), c.getLogger("repository"))
}
