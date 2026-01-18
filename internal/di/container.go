package di

import (
	"context"
	"errors"
	"fmt"

	cacheservice "github.com/michaeldcanady/go-onedrive/internal/app/cache_service"
	clientservice "github.com/michaeldcanady/go-onedrive/internal/app/client_service"
	credentialservice "github.com/michaeldcanady/go-onedrive/internal/app/credential_service"
	driveservice "github.com/michaeldcanady/go-onedrive/internal/app/drive_service"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/event"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	"go.uber.org/zap"
)

type Container struct {
	Ctx                context.Context
	Config             config.Config
	Logger             logging.Logger
	CacheService       CacheService
	CredentialService  CredentialService
	GraphClientService Clienter
	DriveService       ChildrenIterator
	EventBus           *event.InMemoryBus
}

func initializeLogger(logCfg config.LoggingConfig) (logging.Logger, error) {
	if logCfg == nil {
		// TODO: apply default logger config
		return nil, ErrMissingLoggingConfig
	}

	cfg := zap.NewProductionConfig()

	switch logCfg.GetLevel() {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return nil, fmt.Errorf("unknown logging level: %s", logCfg.GetLevel())
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build zap logger: %w", err)
	}

	return logging.NewZapLoggerAdapter(zapLogger), nil
}

func NewContainer(ctx context.Context, cfg config.Config) (*Container, error) {
	var err error

	c := &Container{Ctx: ctx, Config: cfg}

	// logger
	logger, _ := initializeLogger(cfg.GetLoggingConfig())
	c.Logger = logger

	// event bus
	bus := event.NewInMemoryBus(logger)
	c.EventBus = bus

	// services
	if c.CacheService, err = cacheservice.New(cfg.GetAuthenticationConfig().GetProfileCache(), logger); err != nil {
		return nil, errors.Join(errors.New("unable to initialize container"), err)
	}
	c.CredentialService = credentialservice.New(c.CacheService, bus, logger)
	c.GraphClientService = clientservice.New(c.CredentialService, bus, logger)
	c.DriveService = driveservice.New(c.GraphClientService, bus, logger)

	// wiring listeners
	bus.Subscribe(credentialservice.CredentialLoadedTopic,
		event.ListenerFunc(func(ctx context.Context, evt event.Topicer) error {
			_, err := c.GraphClientService.Client(ctx)
			return err
		}),
	)

	return c, nil
}
