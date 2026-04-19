package di

import (
	"github.com/michaeldcanady/go-onedrive/internal/config"
	registry "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/michaeldcanady/go-onedrive/internal/drive"
	"github.com/michaeldcanady/go-onedrive/internal/editor"
	"github.com/michaeldcanady/go-onedrive/internal/environment"
	idregistry "github.com/michaeldcanady/go-onedrive/internal/identity/registry"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/mount"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
)

// Container defines the interface for retrieving and managing core application services.
type Container interface {
	// Logger returns the global logger service.
	Logger() logger.Service
	// Config returns the configuration service.
	Config() config.Service
	// Mounts returns the VFS mount management service.
	Mounts() mount.Service
	// Identity returns the identity provider registry.
	Identity() idregistry.Service
	// Profile returns the configuration profile service.
	Profile() profile.Service
	// FS returns the orchestrated filesystem.
	FS() registry.Service

	// Environment returns the environment service.
	Environment() environment.Service
	// Editor returns the editor service.
	Editor() editor.Service
	// Drive returns the drive-related service.
	Drive() drive.Service
	// URIFactory returns the URI factory service.
	URIFactory() *registry.URIFactory
}
