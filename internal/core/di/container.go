// Package di provides a Dependency Injection container for orchestrating the application's services.
package di

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/core/plugins"
	"github.com/michaeldcanady/go-onedrive/internal/core/resolver"
	"github.com/michaeldcanady/go-onedrive/internal/features/config"
	"github.com/michaeldcanady/go-onedrive/internal/features/drive"
	"github.com/michaeldcanady/go-onedrive/internal/features/editor"
	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	"github.com/michaeldcanady/go-onedrive/internal/features/profile"
	"github.com/michaeldcanady/go-onedrive/internal/features/storage"
	"github.com/michaeldcanady/go-onedrive/internal/features/vfs"
	"github.com/michaeldcanady/go-onedrive/pkg/format"
)

// Container coordinates the access to all shared services and domain logic within the application.
// It serves as the single source of truth for dependency resolution.
type Container interface {
	Logger() logger.Service
	Storage() storage.Service[storage.BoltDB]
	Config() config.Service
	Profile() profile.Service

	PluginManager() plugins.Manager
	VFS() vfs.VFS
	Formatter() format.Factory
	Drive() drive.Service
	Identity() identity.Service
	Token() identity.TokenService
	Mounts() mount.Service
	Editor() editor.Service
	Resolver() resolver.Service

	Shutdown(ctx context.Context) error
}

type container struct {
	logger        logger.Service
	storage       storage.Service[storage.BoltDB]
	config        config.Service
	profile       profile.Service
	pluginManager plugins.Manager
	tokenService  identity.TokenService
	identity      identity.Service
	mounts        mount.Service
	vfs           vfs.VFS
	drive         drive.Service
	editor        editor.Service
	formatter     format.Factory
	resolver      resolver.Service

	services []any
}

// NewContainer returns a new [Container] initialized with the provided service implementations.
func NewContainer(
	l logger.Service,
	s storage.Service[storage.BoltDB],
	c config.Service,
	p profile.Service,
	pm plugins.Manager,
	ts identity.TokenService,
	is identity.Service,
	ms mount.Service,
	v vfs.VFS,
	d drive.Service,
	e editor.Service,
	f format.Factory,
	r resolver.Service,
) Container {
	services := []any{l, s, c, p, pm, ts, is, ms, v, d, e, f, r}

	return &container{
		logger:        l,
		storage:       s,
		config:        c,
		profile:       p,
		pluginManager: pm,
		tokenService:  ts,
		identity:      is,
		mounts:        ms,
		vfs:           v,
		drive:         d,
		editor:        e,
		formatter:     f,
		resolver:      r,
		services:      services,
	}
}

func (c *container) Logger() logger.Service                   { return c.logger }
func (c *container) Storage() storage.Service[storage.BoltDB] { return c.storage }
func (c *container) Config() config.Service                   { return c.config }
func (c *container) Profile() profile.Service                 { return c.profile }

func (c *container) PluginManager() plugins.Manager { return c.pluginManager }
func (c *container) VFS() vfs.VFS                   { return c.vfs }
func (c *container) Formatter() format.Factory      { return c.formatter }
func (c *container) Drive() drive.Service           { return c.drive }
func (c *container) Identity() identity.Service     { return c.identity }
func (c *container) Token() identity.TokenService   { return c.tokenService }
func (c *container) Mounts() mount.Service          { return c.mounts }
func (c *container) Editor() editor.Service         { return c.editor }
func (c *container) Resolver() resolver.Service     { return c.resolver }

func (c *container) Shutdown(ctx context.Context) error {
	var errs []error
	for i := len(c.services) - 1; i >= 0; i-- {
		s := c.services[i]
		if shutdowner, ok := s.(Shutdowner); ok {
			if err := shutdowner.Shutdown(ctx); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown failed with %d errors", len(errs))
	}
	return nil
}
