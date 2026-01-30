package di

import (
	domainauth "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	domainenvironment "github.com/michaeldcanady/go-onedrive/internal2/domain/common/environment"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/config"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	domainprofile "github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
)

type Container interface {
	CacheService() domaincache.CacheService
	FS() domainfs.Service
	EnvironmentService() domainenvironment.EnvironmentService
	Logger() domainlogger.LoggerService
	Auth() domainauth.AuthService
	Profile() domainprofile.ProfileService
	Config() config.ConfigService
	File() file.FileService
}
