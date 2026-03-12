package di

import (
	"github.com/michaeldcanady/go-onedrive/internal2/core/config"
	"github.com/michaeldcanady/go-onedrive/internal2/core/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/core/editor"
	"github.com/michaeldcanady/go-onedrive/internal2/core/environment"
	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/registry"
	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/shared"
	idregistry "github.com/michaeldcanady/go-onedrive/internal2/core/identity/registry"
	"github.com/michaeldcanady/go-onedrive/internal2/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/core/profile"
	"github.com/michaeldcanady/go-onedrive/internal2/core/state"
)

// Container defines the interface for retrieving and managing core application services.
type Container interface {
	// Logger returns the global logger service.
	Logger() logger.Service
	// Config returns the configuration service.
	Config() config.Service
	// State returns the application state service.
	State() state.Service
	// Identity returns the identity provider registry.
	Identity() idregistry.Service
	// Profile returns the configuration profile service.
	Profile() profile.Service
	// FS returns the filesystem provider registry.
	FS() registry.Service
	// Manager returns the orchestrated filesystem manager.
	Manager() shared.Service
	// Environment returns the environment service.
	Environment() environment.Service
	// Editor returns the editor service.
	Editor() editor.Service
	// Drive returns the drive-related service.
	Drive() drive.Service
}
