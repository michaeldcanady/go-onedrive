package di

import (
	"context"
	"errors"
	"path/filepath"

	cacheservice "github.com/michaeldcanady/go-onedrive/internal/app/cache_service"
	clientservice "github.com/michaeldcanady/go-onedrive/internal/app/client_service"
	credentialservice "github.com/michaeldcanady/go-onedrive/internal/app/credential_service"
	driveservice "github.com/michaeldcanady/go-onedrive/internal/app/drive_service"
	environmentservice "github.com/michaeldcanady/go-onedrive/internal/app/environment_service"
	loggerservice "github.com/michaeldcanady/go-onedrive/internal/app/logger_service"
	"github.com/michaeldcanady/go-onedrive/internal/config"
	"github.com/michaeldcanady/go-onedrive/internal/event"
)

type Container struct {
	Ctx                context.Context
	CacheService       CacheService
	CredentialService  CredentialService
	EnvironmentService EnvironmentService
	GraphClientService Clienter
	DriveService       ChildrenIterator
	EventBus           *event.InMemoryBus
	LoggerService      LoggerService
}

func NewContainer(ctx context.Context, _ config.Config) (*Container, error) {
	var err error

	c := &Container{Ctx: ctx}

	// logger
	c.EnvironmentService = environmentservice.New("odc")
	logDir, err := c.EnvironmentService.LogDir(ctx)
	if err != nil {
		return nil, err
	}
	if c.LoggerService, err = loggerservice.New("debug", logDir); err != nil {
		return nil, err
	}

	// event bus
	busLogger, err := c.LoggerService.CreateLogger("bus")
	if err != nil {
		return nil, err
	}
	bus := event.NewInMemoryBus(busLogger)
	c.EventBus = bus

	// services
	cacheLogger, err := c.LoggerService.CreateLogger("cache")
	if err != nil {
		return nil, err
	}

	cacheDir, err := c.EnvironmentService.CacheDir(ctx)
	if err != nil {
		return nil, err
	}

	if c.CacheService, err = cacheservice.New(filepath.Join(cacheDir, "profile.cache"), cacheLogger); err != nil {
		return nil, errors.Join(errors.New("unable to initialize container"), err)
	}
	credentialLogger, err := c.LoggerService.CreateLogger("credential")
	if err != nil {
		return nil, err
	}
	c.CredentialService = credentialservice.New(c.CacheService, bus, credentialLogger)
	graphLogger, err := c.LoggerService.CreateLogger("graph")
	if err != nil {
		return nil, err
	}
	c.GraphClientService = clientservice.New(c.CredentialService, bus, graphLogger)
	driveLogger, err := c.LoggerService.CreateLogger("drive")
	if err != nil {
		return nil, err
	}
	c.DriveService = driveservice.New(c.GraphClientService, bus, driveLogger)

	// wiring listeners
	bus.Subscribe(credentialservice.CredentialLoadedTopic,
		event.ListenerFunc(func(ctx context.Context, evt event.Topicer) error {
			_, err := c.GraphClientService.Client(ctx)
			return err
		}),
	)

	return c, nil
}
