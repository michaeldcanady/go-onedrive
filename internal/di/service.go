package di

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy" // Keep azcore import for policy.TokenRequestOptions
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	registry "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/michaeldcanady/go-onedrive/internal/drive"
	graphgateway "github.com/michaeldcanady/go-onedrive/internal/drive/gateway/graph"
	"github.com/michaeldcanady/go-onedrive/internal/editor"
	"github.com/michaeldcanady/go-onedrive/internal/environment"
	"github.com/michaeldcanady/go-onedrive/internal/identity"
	"github.com/michaeldcanady/go-onedrive/internal/identity/providers/microsoft"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/mount"
	"github.com/michaeldcanady/go-onedrive/internal/profile"
	"github.com/michaeldcanady/go-onedrive/internal/storage/backend/local"
	"github.com/michaeldcanady/go-onedrive/internal/storage/backend/onedrive"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/michaeldcanady/go-onedrive/pkg/logger/zap"
	bolt "go.etcd.io/bbolt"
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
	// identityService is the registry for managing multiple identity providers, providing both authenticators and authorizers.
	identityService identity.Service
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

	profileSvc, err := profile.NewDefaultService(c.environment)
	if err != nil {
		return fmt.Errorf("failed to initialize profile service: %w", err)
	}
	c.profile = profileSvc

	return nil
}

func (c *DefaultContainer) initIdentityServices(ctx context.Context) error {
	cliLog, _ := c.logger.CreateLogger("cli")

	stateDir, err := c.environment.StateDir()
	if err != nil {
		return fmt.Errorf("failed to get state directory: %w", err)
	}
	dbPath := filepath.Join(stateDir, "identity.db")
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("failed to open identity database: %w", err)
	}

	tokenRepo := identity.NewBoltRepository(db)
	// Create the MicrosoftAuthenticator.
	msAuth := microsoft.NewMicrosoftAuthenticator()

	// Create a registry for Identity providers.
	authRegistry := identity.NewRegistry(tokenRepo, cliLog)
	authRegistry.RegisterAuthenticator("microsoft", msAuth) // Register Authenticator

	// Create and register the MicrosoftAuthorizer.
	cred, err := msAuth.GetCredentialForAuthorizer(ctx)
	if err != nil {
		return fmt.Errorf("failed to get token credential for authorizer: %w", err)
	}
	msAuthorizer := microsoft.NewMicrosoftAuthorizer(cred)
	authRegistry.RegisterAuthorizer("microsoft", msAuthorizer)

	c.identityService = authRegistry // Assigning the registry to identityService

	return nil
}

func (c *DefaultContainer) initDriveServices(ctx context.Context) error {
	cliLog, _ := c.logger.CreateLogger("cli")

	// Get the Authorizer for Microsoft.
	msAuthorizer, err := c.identityService.GetAuthorizer("microsoft") // Use GetAuthorizer
	if err != nil {
		return fmt.Errorf("failed to get authorizer for microsoft: %w", err)
	}

	// Pass the Authorizer to the GraphDriveGateway.
	driveGateway := graphgateway.NewGraphDriveGateway(msAuthorizer, cliLog)
	c.drive = drive.NewDefaultService(driveGateway, cliLog)

	return nil
}

func (c *DefaultContainer) initVFSServices(ctx context.Context) error {
	cliLog, _ := c.logger.CreateLogger("cli")
	yamlSvc := config.NewConfigService(c.profile, cliLog)
	c.config = yamlSvc
	c.mounts = mount.NewMountService(c.config)

	vfs := registry.NewVFS()
	appConfig, _ := c.config.GetConfig(ctx)

	// Get the Authorizer for Microsoft.
	msAuthorizer, err := c.identityService.GetAuthorizer("microsoft") // Use GetAuthorizer
	if err != nil {
		return fmt.Errorf("failed to get authorizer for microsoft: %w", err)
	}

	// Helper to create onedrive backend
	createOneDriveBackend := func(identityID, driveID string) pkgfs.Backend {
		// The Authorizer returns an AccessToken. We need azcore.TokenCredential for NewGraphProvider.
		// We will create a StaticTokenCredential from the AccessToken.
		accessToken, tokenErr := msAuthorizer.Token(ctx, identityID)
		if tokenErr != nil {
			cliLog.Error("failed to get token for onedrive backend", logger.Error(tokenErr), logger.String("identity", identityID))
			return nil // Return nil if token acquisition fails.
		}

		// Convert identity.AccessToken to azcore.TokenCredential.
		// NOTE: This relies on the structure of identity.AccessToken and azcore.TokenCredential.
		// This conversion might need further refinement if more complex credential types are used.
		cred := azidentity.NewStaticTokenCredential(accessToken.Token, accessToken.ExpiresAt, accessToken.Scopes)

		p := microsoft.NewGraphProvider(cred, cliLog)
		dr := drive.NewDefaultResolver(c.drive, identityID)
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
	c.uriFactory = registry.NewURIFactory(vfs)

	cliLog, _ = c.logger.CreateLogger("cli")
	c.editor = editor.NewDefaultService(c.environment, c.uriFactory, cliLog, editor.WithConfig(c.config))

	return nil
}

// Logger returns the global logging service.
func (c *DefaultContainer) Logger() logger.Service { return c.logger }

// Config returns the configuration management service.
func (c *DefaultContainer) Config() config.Service { return c.config }

// Mounts returns the VFS mount management service.
func (c *DefaultContainer) Mounts() mount.Service { return c.mounts }

// Identity returns the identity provider registry.
func (c *DefaultContainer) Identity() identity.Service { return c.identityService }

// Profile returns the configuration profile service.
func (c *DefaultContainer) Profile() profile.Service { return c.profile }

// FS returns the orchestrated filesystem.
func (c *DefaultContainer) FS() registry.Service { return c.manager }

// Environment returns the environment-related service.
func (c *DefaultContainer) Environment() environment.Service { return c.environment }

// Editor returns the external editor service.
func (c *DefaultContainer) Editor() editor.Service { return c.editor }

// Drive returns the OneDrive drive management service.
func (c *DefaultContainer) Drive() drive.Service { return c.drive }

// URIFactory returns the URI factory service.
func (c *DefaultContainer) URIFactory() *registry.URIFactory {
	return c.uriFactory
}
 Here is the updated code:
...
	tokenRepo := identity.NewBoltRepository(db)
	// Create the MicrosoftAuthenticator.
	msAuth := microsoft.NewMicrosoftAuthenticator(tokenRepo, cliLog)

	// Create a registry for Identity providers.
	authRegistry := identity.NewRegistry(cliLog) // Pass logger to NewRegistry
	authRegistry.RegisterAuthenticator("microsoft", msAuth) // Register Authenticator

	// Create and register the MicrosoftAuthorizer.
	// This requires an azcore.TokenCredential.
	// For DI, this credential creation might be complex and depend on config/state.
	// For now, let's assume a way to get a credential is provided or can be created.
	// This part needs careful DI wiring. As a placeholder, we might need to create a minimal credential.
	// A better approach would be for the DI system to provide this directly.
	// For now, let's assume `getOrUpdateCredential` from MicrosoftAuthenticator can provide it,
	// although that couples the DI too tightly with the implementation detail.
	// A more appropriate place to create this credential might be in the DI setup itself.
	// For now, we'll create a placeholder.
	// TODO: Refine credential creation for authorizer.

	// Create the MicrosoftAuthorizer. This requires an azcore.TokenCredential.
	// We need to get this credential. For now, we'll reuse the logic from MicrosoftAuthenticator.
	// This is a placeholder and needs proper DI wiring for credential management.
	// In a production scenario, the credential itself (or a factory for it) would be managed by DI.
	msAuthCred, err := msAuth.GetCredentialForAuthorizer(ctx) // Assuming this method exists or will be added to MicrosoftAuthenticator
	if err != nil {
		return nil, fmt.Errorf("failed to get token credential for authorizer: %w", err)
	}
	msAuthorizer := microsoft.NewMicrosoftAuthorizer(msAuthCred, tokenRepo, cliLog)
	authRegistry.RegisterAuthorizer("microsoft", msAuthorizer) // Register Authorizer

	c.identityService = authRegistry // Assigning the registry to identityService

	return nil
}

func (c *DefaultContainer) initDriveServices(ctx context.Context) error {
	cliLog, _ := c.logger.CreateLogger("cli")

	// Get the Authorizer for Microsoft.
	// This assumes that the identity.Service (Registry) has a method to retrieve an Authorizer.
	msAuthorizer, err := c.identityService.GetAuthorizer("microsoft") // Use GetAuthorizer
	if err != nil {
		return fmt.Errorf("failed to get authorizer for microsoft: %w", err)
	}

	// Pass the Authorizer to the GraphDriveGateway.
	driveGateway := graphgateway.NewGraphDriveGateway(msAuthorizer, cliLog)
	c.drive = drive.NewDefaultService(driveGateway, cliLog)

	return nil
}

func (c *DefaultContainer) initVFSServices(ctx context.Context) error {
	cliLog, _ := c.logger.CreateLogger("cli")
	yamlSvc := config.NewConfigService(c.profile, cliLog)
	c.config = yamlSvc
	c.mounts = mount.NewMountService(c.config)

	vfs := registry.NewVFS()
	appConfig, _ := c.config.GetConfig(ctx)

	// Get the Authorizer for Microsoft.
	msAuthorizer, err := c.identityService.GetAuthorizer("microsoft") // Use GetAuthorizer
	if err != nil {
		return fmt.Errorf("failed to get authorizer for microsoft: %w", err)
	}

	// Helper to create onedrive backend
	createOneDriveBackend := func(identityID, driveID string) pkgfs.Backend {
		// The Authorizer returns an AccessToken. We need azcore.TokenCredential for NewGraphProvider.
		// We will create a StaticTokenCredential from the AccessToken.
		accessToken, tokenErr := msAuthorizer.Token(ctx, identityID)
		if tokenErr != nil {
			cliLog.Error("failed to get token for onedrive backend", logger.Error(tokenErr), logger.String("identity", identityID))
			return nil // Return nil if token acquisition fails.
		}

		// Convert identity.AccessToken to azcore.TokenCredential.
		// NOTE: This relies on the structure of identity.AccessToken and azcore.TokenCredential.
		// This conversion might need further refinement if more complex credential types are used.
		cred := azidentity.NewStaticTokenCredential(accessToken.Token, accessToken.ExpiresAt, accessToken.Scopes)

		p := microsoft.NewGraphProvider(cred, cliLog)
		dr := drive.NewDefaultResolver(c.drive, identityID)
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
	c.uriFactory = registry.NewURIFactory(vfs)

	cliLog, _ = c.logger.CreateLogger("cli")
	c.editor = editor.NewDefaultService(c.environment, c.uriFactory, cliLog, editor.WithConfig(c.config))

	return nil
}

// Logger returns the global logging service.
func (c *DefaultContainer) Logger() logger.Service { return c.logger }

// Config returns the configuration management service.
func (c *DefaultContainer) Config() config.Service { return c.config }

// Mounts returns the VFS mount management service.
func (c *DefaultContainer) Mounts() mount.Service { return c.mounts }

// Identity returns the identity provider registry.
func (c *DefaultContainer) Identity() identity.Service { return c.identityService }

// Profile returns the configuration profile service.
func (c *DefaultContainer) Profile() profile.Service { return c.profile }

// FS returns the orchestrated filesystem.
func (c *DefaultContainer) FS() registry.Service { return c.manager }

// Environment returns the environment-related service.
func (c *DefaultContainer) Environment() environment.Service { return c.environment }

// Editor returns the external editor service.
func (c *DefaultContainer) Editor() editor.Service { return c.editor }

// Drive returns the OneDrive drive management service.
func (c *DefaultContainer) Drive() drive.Service { return c.drive }

// URIFactory returns the URI factory service.
func (c *DefaultContainer) URIFactory() *registry.URIFactory {
	return c.uriFactory
}
