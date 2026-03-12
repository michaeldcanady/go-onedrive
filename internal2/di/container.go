package di

import (
	"github.com/michaeldcanady/go-onedrive/internal2/core/config"
	"github.com/michaeldcanady/go-onedrive/internal2/core/identity/registry"
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
	Identity() registry.Service
	// Profile returns the configuration profile service.
	Profile() profile.Service
}
