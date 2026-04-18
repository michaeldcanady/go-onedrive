package di

import (
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/drive"
	"github.com/michaeldcanady/go-onedrive/internal/alias"
	"github.com/michaeldcanady/go-onedrive/internal/environment"
	registry "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/michaeldcanady/go-onedrive/internal/editor"
	idregistry "github.com/michaeldcanady/go-onedrive/internal/identity/registry"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/mount"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
	"github.com/michaeldcanady/go-onedrive/internal/state"
)

// Container defines the interface for retrieving and managing core application services.
type Container interface {
	// Logger returns the global logger service.
	Logger() logger.Service
	// Config returns the configuration service.
	Config() config.Service
	// Mounts returns the VFS mount management service.
	Mounts() mount.Service
	// State returns the application state service.
	State() state.Service
	// Identity returns the identity provider registry.
	Identity() idregistry.Service
	// Profile returns the configuration profile service.
	Profile() profile.Service
	// FS returns the orchestrated filesystem.
	FS() registry.Service

	ProviderRegistry() interface {
		RegisteredNames() ([]string, error)
	}
	// Environment returns the environment service.
	Environment() environment.Service
	// Editor returns the editor service.
	Editor() editor.Service
	// Drive returns the drive-related service.
	Drive() drive.Service
	// Alias returns the drive alias management service.
	Alias() alias.Service
	// URIFactory returns the URI factory service.
	URIFactory() *registry.URIFactory
}
