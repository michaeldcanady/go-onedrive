package di

import (
	"github.com/michaeldcanady/go-onedrive/internal2/core/config"
	"github.com/michaeldcanady/go-onedrive/internal2/core/identity/providers/microsoft"
	"github.com/michaeldcanady/go-onedrive/internal2/core/identity/registry"
	"github.com/michaeldcanady/go-onedrive/internal2/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/core/profile"
	"github.com/michaeldcanady/go-onedrive/internal2/core/state"
)

// DefaultContainer provides a concrete implementation of the Container interface.
// It orchestrates the lifecycle and wiring of core application services.
type DefaultContainer struct {
	// logger is the centralized logging service.
	logger logger.Service
	// config is the service for managing application settings.
	config config.Service
	// state is the service for tracking session and persistent state.
	state state.Service
	// identity is the registry for managing multiple identity providers.
	identity registry.Service
	// profile is the service for managing user configuration profiles.
	profile profile.Service
}

// NewDefaultContainer initializes a new instance of the DefaultContainer with all core services wired.
func NewDefaultContainer() *DefaultContainer {
	logSvc := logger.NewZapService()
	cliLog, _ := logSvc.CreateLogger("cli")

	idReg := registry.NewRegistry()
	idReg.Register("microsoft", microsoft.NewAuthenticator(nil, cliLog))

	return &DefaultContainer{
		logger:   logSvc,
		config:   config.NewYAMLService(cliLog),
		state:    state.NewMemoryService(),
		identity: idReg,
		profile:  profile.NewMemoryService(),
	}
}

// Logger returns the global logging service.
func (c *DefaultContainer) Logger() logger.Service { return c.logger }

// Config returns the configuration management service.
func (c *DefaultContainer) Config() config.Service { return c.config }

// State returns the application state tracking service.
func (c *DefaultContainer) State() state.Service { return c.state }

// Identity returns the identity provider registry.
func (c *DefaultContainer) Identity() registry.Service { return c.identity }

// Profile returns the configuration profile service.
func (c *DefaultContainer) Profile() profile.Service { return c.profile }
