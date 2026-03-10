package di

import (
	"github.com/michaeldcanady/go-onedrive/internal2/app/common/environment"
	applogging "github.com/michaeldcanady/go-onedrive/internal2/app/common/logging"
	appconfig "github.com/michaeldcanady/go-onedrive/internal2/app/config"
	domainenv "github.com/michaeldcanady/go-onedrive/internal2/domain/common/environment"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	domainconfig "github.com/michaeldcanady/go-onedrive/internal2/domain/config"
	infralogging "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	infraconfig "github.com/michaeldcanady/go-onedrive/internal2/infra/config"
)

func (c *Container) getLogger(name string) domainlogger.Logger {
	loggerService := c.Logger()
	log, err := loggerService.CreateLogger(name)
	if err != nil {
		return infralogging.NewNoopLogger()
	}
	return log
}

// EnvironmentService implements [di.Container].
func (c *Container) EnvironmentService() domainenv.EnvironmentService {
	c.environmentOnce.Do(func() {
		c.environmentService = c.newEnvironmentService()
	})
	return c.environmentService
}

func (c *Container) newEnvironmentService() domainenv.EnvironmentService {
	svc := environment.New2("odc")
	_ = svc.EnsureAll()
	return svc
}

// Logger implements [di.Container].
func (c *Container) Logger() domainlogger.LoggerService {
	c.loggerOnce.Do(func() {
		c.loggerService = c.newLoggerService()
	})
	return c.loggerService
}

func (c *Container) newLoggerService() domainlogger.LoggerService {
	level, _ := c.EnvironmentService().LogLevel()

	opts := []domainlogger.Option{domainlogger.WithLogLevel(level), domainlogger.WithType(domainlogger.TypeZap)}

	outputDest, _ := c.EnvironmentService().OutputDestination()
	switch outputDest {
	case domainlogger.OutputDestinationFile:
		logHome, _ := c.EnvironmentService().LogDir()
		opts = append(opts, domainlogger.WithOutputDestinationFile(logHome))
	case domainlogger.OutputDestinationStandardOut:
		opts = append(opts, domainlogger.WithOutputDestinationStandardOut())
	case domainlogger.OutputDestinationStandardError:
		opts = append(opts, domainlogger.WithOutputDestinationStandardError())
	default:
	}

	svc, _ := applogging.NewLoggerService(opts...)
	svc.RegisterProvider(domainlogger.TypeZap, infralogging.NewZapLoggerProvider())
	return svc
}

// Config implements [di.Container].
func (c *Container) Config() domainconfig.ConfigService {
	c.configOnce.Do(func() {
		c.configService = c.newConfigService()
	})
	return c.configService
}

func (c *Container) newConfigService() domainconfig.ConfigService {
	return appconfig.New2(infraconfig.NewYAMLLoader(), c.getLogger("config"))
}
