package di

import (
	"github.com/michaeldcanady/go-onedrive/internal/feature/config"
	"github.com/michaeldcanady/go-onedrive/internal/feature/drive"
	"github.com/michaeldcanady/go-onedrive/internal/feature/environment"
	registry "github.com/michaeldcanady/go-onedrive/internal/feature/fs"
	"github.com/michaeldcanady/go-onedrive/internal/feature/fs/editor"
	idregistry "github.com/michaeldcanady/go-onedrive/internal/feature/identity/registry"
	"github.com/michaeldcanady/go-onedrive/internal/feature/logger"
	"github.com/michaeldcanady/go-onedrive/internal/feature/profile"
	"github.com/michaeldcanady/go-onedrive/internal/feature/state"
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
	// FS returns the orchestrated filesystem.
	FS() registry.Service
	// Environment returns the environment service.
	Environment() environment.Service
	// Editor returns the editor service.
	Editor() editor.Service
	// Drive returns the drive-related service.
	Drive() drive.Service
}
