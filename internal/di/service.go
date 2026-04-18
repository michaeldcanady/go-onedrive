package di

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/drive"
	"github.com/michaeldcanady/go-onedrive/internal/alias"
	graphgateway "github.com/michaeldcanady/go-onedrive/internal/drive/gateway/graph"
	"github.com/michaeldcanady/go-onedrive/internal/environment"
	registry "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/michaeldcanady/go-onedrive/internal/storage/backend/local"
	"github.com/michaeldcanady/go-onedrive/internal/storage/backend/onedrive"
	"github.com/michaeldcanady/go-onedrive/internal/editor"
	"github.com/michaeldcanady/go-onedrive/internal/identity/providers/microsoft"
	idregistry "github.com/michaeldcanady/go-onedrive/internal/identity/registry"
	idshared "github.com/michaeldcanady/go-onedrive/internal/identity/shared"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/mount"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
	"github.com/michaeldcanady/go-onedrive/internal/state"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/logger/zap"
)

// DefaultContainer provides a concrete implementation of the Container interface.
// It orchestrates the lifecycle and wiring of core application services.
type DefaultContainer struct {
	// logger is the centralized logging service.
	logger logger.Service
	// config is the service for managing application settings.
	config config.Service
	// mounts is the service for managing VFS mount points.
	mounts mount.Service
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
}

// NewDefaultContainer initializes a new instance of the DefaultContainer with all core services wired.
func NewDefaultContainer() (*DefaultContainer, error) {
	ctx := context.Background()
	c := &DefaultContainer{}

	if err := c.initBaseServices(); err != nil {
		return nil, err
	}

	if err := c.initIdentityServices(ctx); err != nil {
		return nil, err
	}

	if err := c.initDriveServices(ctx); err != nil {
		return nil, err
	}

	if err := c.initVFSServices(ctx); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *DefaultContainer) initBaseServices() error {
	c.environment = environment.NewDefaultService("odc")
	if err := c.environment.EnsureAll(); err != nil {
		return fmt.Errorf("failed to ensure environment directories: %w", err)
	}

	c.logger = zap.NewZapService(c.environment)

	stateSvc, err := state.NewBoltService(c.environment)
	if err != nil {
		return fmt.Errorf("failed to initialize state service: %w", err)
	}
	c.state = stateSvc

	profileSvc, err := profile.NewDefaultService(c.environment, c.state)
	if err != nil {
		return fmt.Errorf("failed to initialize profile service: %w", err)
	}
	c.profile = profileSvc

	cliLog, _ := c.logger.CreateLogger("cli")
	c.editor = editor.NewDefaultService(c.environment, cliLog)

	return nil
}

func (c *DefaultContainer) initIdentityServices(ctx context.Context) error {
	cliLog, _ := c.logger.CreateLogger("cli")
	msAuth := microsoft.NewAuthenticator(c.state, cliLog)

	c.identity = idregistry.NewRegistry()
	c.identity.Register("microsoft", msAuth)

	// Legacy token support
	tokenData, err := c.state.Get(state.KeyAccessToken)
	if err == nil && tokenData != "" {
		var token idshared.AccessToken
		_ = json.Unmarshal([]byte(tokenData), &token)
	}

	return nil
}

func (c *DefaultContainer) initDriveServices(ctx context.Context) error {
	cliLog, _ := c.logger.CreateLogger("cli")
	msAuth, _ := c.identity.Get("microsoft")

	driveGateway := graphgateway.NewGraphDriveGateway(msAuth, cliLog)
	c.drive = drive.NewDefaultService(driveGateway, c.state, cliLog)

	aliasSvc, err := alias.NewDefaultService(c.environment, cliLog)
	if err != nil {
		return fmt.Errorf("failed to initialize drive alias service: %w", err)
	}
	c.alias = aliasSvc

	return nil
}

func (c *DefaultContainer) initVFSServices(ctx context.Context) error {
	cliLog, _ := c.logger.CreateLogger("cli")
	yamlSvc := config.NewConfigService(c.profile, c.state, cliLog)
	c.config = yamlSvc
	c.mounts = mount.NewMountService(c.config)

	vfs := registry.NewVFS()
	appConfig, _ := c.config.GetConfig(ctx)

	msAuth, _ := c.identity.Get("microsoft")

	// Helper to create onedrive backend
	createOneDriveBackend := func(identityID, driveID string) pkgfs.Backend {
		var cachedCred azcore.TokenCredential
		if cred, err := msAuth.GetCredential(ctx, identityID); err == nil {
			cachedCred = cred.(azcore.TokenCredential)
		}
		p := microsoft.NewGraphProvider(cachedCred, cliLog)
		dr := drive.NewDefaultResolver(c.state, identityID)
		return onedrive.NewBackend(p, driveID, dr, cliLog)
	}

	if len(appConfig.Mounts) == 0 {
		localBackend := local.NewBackend("/", cliLog)
		onedriveBackend := createOneDriveBackend("", "")

		vfs.Mount("/", localBackend)
		vfs.Mount("/local", localBackend)
		vfs.Mount("/onedrive", onedriveBackend)
	} else {
		for _, m := range appConfig.Mounts {
			var backend pkgfs.Backend
			switch m.Type {
			case "local":
				root := m.Options["root"]
				if root == "" {
					root = "/"
				}
				backend = local.NewBackend(root, cliLog)
			case "onedrive":
				backend = createOneDriveBackend(m.IdentityID, m.Options["drive_id"])
			default:
				cliLog.Warn("unknown backend type in config", logger.String("type", m.Type), logger.String("path", m.Path))
				continue
			}
			if backend != nil {
				vfs.Mount(m.Path, backend)
			}
		}
	}

	c.manager = vfs
	c.uriFactory = registry.NewURIFactory(vfs, c.alias)

	return nil
}

// Logger returns the global logging service.
func (c *DefaultContainer) Logger() logger.Service { return c.logger }

// Config returns the configuration management service.
func (c *DefaultContainer) Config() config.Service { return c.config }

// Mounts returns the VFS mount management service.
func (c *DefaultContainer) Mounts() mount.Service { return c.mounts }

// State returns the application state tracking service.
func (c *DefaultContainer) State() state.Service { return c.state }

// Identity returns the identity provider registry.
func (c *DefaultContainer) Identity() idregistry.Service { return c.identity }

// Profile returns the configuration profile service.
func (c *DefaultContainer) Profile() profile.Service { return c.profile }

func (c *DefaultContainer) ProviderRegistry() interface {
	RegisteredNames() ([]string, error)
} {
	return nil // Obsolete
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
