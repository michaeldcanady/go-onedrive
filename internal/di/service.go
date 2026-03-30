package di

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/michaeldcanady/go-onedrive/internal/core/config"
	"github.com/michaeldcanady/go-onedrive/internal/core/drive"
	"github.com/michaeldcanady/go-onedrive/internal/core/editor"
	"github.com/michaeldcanady/go-onedrive/internal/core/environment"
	"github.com/michaeldcanady/go-onedrive/internal/core/fs/manager"
	"github.com/michaeldcanady/go-onedrive/internal/core/fs/providers/local"
	"github.com/michaeldcanady/go-onedrive/internal/core/fs/providers/onedrive"
	"github.com/michaeldcanady/go-onedrive/internal/core/fs/registry"
	"github.com/michaeldcanady/go-onedrive/internal/core/fs/shared"
	"github.com/michaeldcanady/go-onedrive/internal/core/identity/providers/microsoft"
	idregistry "github.com/michaeldcanady/go-onedrive/internal/core/identity/registry"
	idshared "github.com/michaeldcanady/go-onedrive/internal/core/identity/shared"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/core/profile"
	corems "github.com/michaeldcanady/go-onedrive/internal/core/providers/microsoft"
	"github.com/michaeldcanady/go-onedrive/internal/core/state"
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

	envSvc := environment.NewDefaultService("odc")

	stateSvc, err := state.NewBoltService(envSvc)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize state service: %w", err)
	}
	c.state = stateSvc

	profileSvc, err := profile.NewBoltService(envSvc) // Use BoltService for profiles
	if err != nil {
		return nil, fmt.Errorf("failed to initialize profile service: %w", err)
	}
	c.profile = profileSvc

	// Try to load cached token
	var cachedCred azcore.TokenCredential
	tokenData, err := stateSvc.Get(state.KeyAccessToken)
	if err == nil && tokenData != "" {
		var token idshared.AccessToken
		if err := json.Unmarshal([]byte(tokenData), &token); err == nil {
			cachedCred = microsoft.NewStaticTokenCredential(token)
		} else {
			cliLog.Warn("failed to unmarshal cached token", logger.Error(err))
		}
	}

	msAuth := microsoft.NewAuthenticator(cachedCred, cliLog)
	idReg := idregistry.NewRegistry()
	idReg.Register("microsoft", msAuth)

	graphProvider := corems.NewGraphProvider(msAuth.Credential(), cliLog)

	driveGateway := onedrive.NewGraphDriveGateway(graphProvider, cliLog)
	driveSvc := drive.NewDefaultService(driveGateway, cliLog)

	fsReg := registry.NewRegistry(stateSvc)
	fsReg.Register("local", local.NewProvider(cliLog))
	fsReg.Register("onedrive", onedrive.NewProvider(graphProvider, stateSvc, driveSvc, cliLog))

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
