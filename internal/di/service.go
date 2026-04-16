package di

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/drive"
	"github.com/michaeldcanady/go-onedrive/internal/drive/alias"
	graphgateway "github.com/michaeldcanady/go-onedrive/internal/drive/gateway/graph"
	"github.com/michaeldcanady/go-onedrive/internal/environment"
	registry "github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/fs/editor"
	"github.com/michaeldcanady/go-onedrive/internal/fs/providers"
	_ "github.com/michaeldcanady/go-onedrive/internal/fs/providers/all"
	"github.com/michaeldcanady/go-onedrive/internal/identity/providers/microsoft"
	idregistry "github.com/michaeldcanady/go-onedrive/internal/identity/registry"
	idshared "github.com/michaeldcanady/go-onedrive/internal/identity/shared"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
	"github.com/michaeldcanady/go-onedrive/internal/state"
	"github.com/michaeldcanady/go-onedrive/pkg/logger/zap"
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
	// manager is the orchestrated filesystem manager.
	manager registry.Service
	// environment is the environment-related service.
	environment environment.Service
	// editor is the external editor service.
	editor editor.Service
	// drive is the OneDrive drive management service.
	drive drive.Service

	// alias is the drive alias management service.
	alias alias.Service

	// uriFactory is the service for creating and resolving URIs.
	uriFactory *registry.URIFactory

	registry interface {
		RegisteredNames() ([]string, error)
	}
}

// NewDefaultContainer initializes a new instance of the DefaultContainer with all core services wired.
func NewDefaultContainer() (*DefaultContainer, error) {
	envSvc := environment.NewDefaultService("odc")

	if err := envSvc.EnsureAll(); err != nil {
		return nil, fmt.Errorf("failed to ensure environment directories: %w", err)
	}

	logSvc := zap.NewZapService(envSvc)
	cliLog, _ := logSvc.CreateLogger("cli")

	stateSvc, err := state.NewBoltService(envSvc)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize state service: %w", err)
	}

	profileSvc, err := profile.NewBoltService(envSvc, stateSvc)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize profile service: %w", err)
	}

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

	msAuth := microsoft.NewAuthenticator(cachedCred, stateSvc, cliLog)
	idReg := idregistry.NewRegistry()
	idReg.Register("microsoft", msAuth)

	graphProvider := microsoft.NewGraphProvider(msAuth.Credential(), cliLog)

	driveGateway := graphgateway.NewGraphDriveGateway(graphProvider, cliLog)
	driveSvc := drive.NewDefaultService(driveGateway, stateSvc, cliLog)
	aliasSvc, err := alias.NewBoltService(envSvc, cliLog)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize drive alias service: %w", err)
	}

	fsReg := registry.NewRegistry(stateSvc, aliasSvc, cliLog)

	deps := &providerDeps{
		logger: cliLog,
		values: map[string]any{
			"platform":       graphProvider,
			"drive_resolver": &driveResolver{state: stateSvc},
		},
	}

	for _, name := range providers.RegisteredNames() {
		desc, _ := providers.Get(name)
		svc, err := desc.Factory(deps)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize provider %s: %w", name, err)
		}
		fsReg.Register(name, svc)
	}

	editorSvc := editor.NewDefaultService(envSvc, cliLog)

	configSvc := config.NewYAMLService(profileSvc, stateSvc, cliLog)

	uriFactory := registry.NewURIFactory(fsReg, aliasSvc)

	return &DefaultContainer{
		logger:      logSvc,
		config:      configSvc,
		state:       stateSvc,
		identity:    idReg,
		profile:     profileSvc,
		registry:    fsReg,
		manager:     registry.NewFileSystemManager(fsReg),
		environment: envSvc,
		editor:      editorSvc,
		alias:       aliasSvc,
		drive:       driveSvc,
		uriFactory:  uriFactory,
	}, nil
}

// driveResolver implements fs.DriveResolver using the internal state service.
type driveResolver struct {
	state state.Service
}

func (r *driveResolver) GetActiveDriveID(ctx context.Context) (string, error) {
	return r.state.Get(state.KeyDrive)
}

type providerDeps struct {
	logger logger.Logger
	values map[string]any
}

func (d *providerDeps) Logger() logger.Logger {
	return d.logger
}

func (d *providerDeps) Get(key string) (any, bool) {
	val, ok := d.values[key]
	return val, ok
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

func (c *DefaultContainer) ProviderRegistry() interface {
	RegisteredNames() ([]string, error)
} {
	return c.registry
}

// FS returns the orchestrated filesystem.
func (c *DefaultContainer) FS() registry.Service { return c.manager }

// Environment returns the environment-related service.
func (c *DefaultContainer) Environment() environment.Service { return c.environment }

// Editor returns the external editor service.
func (c *DefaultContainer) Editor() editor.Service { return c.editor }

// Drive returns the OneDrive drive management service.
func (c *DefaultContainer) Drive() drive.Service { return c.drive }

func (c *DefaultContainer) Alias() alias.Service {
	return c.alias
}

// URIFactory returns the URI factory service.
func (c *DefaultContainer) URIFactory() *registry.URIFactory {
	return c.uriFactory
}
