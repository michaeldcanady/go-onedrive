package di

import (
	"context"

	infrafile "github.com/michaeldcanady/go-onedrive/internal2/infra/file"
)

func (c *Container) metadata() *infrafile.MetadataRepository {
	c.metadataOnce.Do(func() {
		c.metadataRepo = c.newMetadataRepository()
	})
	return c.metadataRepo
}

func (c *Container) newMetadataRepository() *infrafile.MetadataRepository {
	adapter, _ := c.clientProvider().RequestAdapter(context.Background())

	return infrafile.NewMetadataRepository(adapter, c.metadataCache(), c.metadataListingCache(), c.pathIDCache(), c.getLogger("repository"))
}

func (c *Container) contents() *infrafile.ContentsRepository {
	c.contentsOnce.Do(func() {
		c.contentsRepo = c.newContentsRepository()
	})
	return c.contentsRepo
}

func (c *Container) newContentsRepository() *infrafile.ContentsRepository {
	adapter, _ := c.clientProvider().RequestAdapter(context.Background())

	return infrafile.NewContentsRepository(adapter, c.contentsCache(), c.metadataCache(), c.pathIDCache(), c.getLogger("repository"))
}
