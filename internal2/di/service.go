package di

import (
	"github.com/michaeldcanady/go-onedrive/internal2/core/config"
	"github.com/michaeldcanady/go-onedrive/internal2/core/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/core/editor"
	"github.com/michaeldcanady/go-onedrive/internal2/core/environment"
	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/manager"
	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/providers/local"
	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/providers/onedrive"
	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/registry"
	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/shared"
	"github.com/michaeldcanady/go-onedrive/internal2/core/identity/providers/microsoft"
	idregistry "github.com/michaeldcanady/go-onedrive/internal2/core/identity/registry"
	"github.com/michaeldcanady/go-onedrive/internal2/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/core/profile"
	corems "github.com/michaeldcanady/go-onedrive/internal2/core/providers/microsoft"
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
	identity idregistry.Service
	// profile is the service for managing user configuration profiles.
	profile profile.Service
	// fs is the registry for filesystem providers.
	fs registry.Service
	// manager is the orchestrated filesystem manager.
	manager shared.Service
	// environment is the environment-related service.
	environment environment.Service
	// editor is the external editor service.
	editor editor.Service
	// drive is the OneDrive drive management service.
	drive drive.Service
}

// NewDefaultContainer initializes a new instance of the DefaultContainer with all core services wired.
func NewDefaultContainer() (*DefaultContainer, error) {
	logSvc := logger.NewZapService()
	cliLog, _ := logSvc.CreateLogger("cli")

	stateSvc, err := state.NewFileService("")
	if err != nil {
		return nil, err
	}

	profileSvc, err := profile.NewFileService("")
	if err != nil {
		return nil, err
	}

	// For identity, we need to decide how to load existing credentials.
	// For now, initializing with nil and allowing login to populate it.
	msAuth := microsoft.NewAuthenticator(nil, cliLog)
	idReg := idregistry.NewRegistry()
	idReg.Register("microsoft", msAuth)

	graphProvider := corems.NewGraphProvider(msAuth.Credential(), cliLog)

	driveGateway := onedrive.NewGraphDriveGateway(graphProvider, cliLog)
	driveSvc := drive.NewDefaultService(driveGateway, cliLog)

	fsReg := registry.NewRegistry(stateSvc)
	fsReg.Register("local", local.NewProvider(cliLog))
	fsReg.Register("onedrive", onedrive.NewProvider(graphProvider, stateSvc, driveSvc, cliLog))

	envSvc := environment.NewDefaultService("odc")
	editorSvc := editor.NewDefaultService(envSvc, cliLog)

	return &DefaultContainer{
		logger:      logSvc,
		config:      config.NewYAMLService(cliLog),
		state:       stateSvc,
		identity:    idReg,
		profile:     profileSvc,
		fs:          fsReg,
		manager:     manager.NewFileSystemManager(fsReg),
		environment: envSvc,
		editor:      editorSvc,
		drive:       driveSvc,
	}, nil
}

// Logger returns the global logging service.
func (c *DefaultContainer) Logger() logger.Service { return c.logger }

// Config returns the configuration management service.
func (c *DefaultContainer) Config() config.Service { return c.config }

// State returns the application state tracking service.
func (c *DefaultContainer) State() state.Service { return c.state }

// Identity returns the identity provider registry.
func (c *DefaultContainer) Identity() idregistry.Service { return c.identity }

// Profile returns the configuration profile service.
func (c *DefaultContainer) Profile() profile.Service { return c.profile }

// FS returns the filesystem provider registry.
func (c *DefaultContainer) FS() registry.Service { return c.fs }

// Manager returns the orchestrated filesystem manager.
func (c *DefaultContainer) Manager() shared.Service { return c.manager }

// Environment returns the environment-related service.
func (c *DefaultContainer) Environment() environment.Service { return c.environment }

// Editor returns the external editor service.
func (c *DefaultContainer) Editor() editor.Service { return c.editor }

// Drive returns the OneDrive drive management service.
func (c *DefaultContainer) Drive() drive.Service { return c.drive }
