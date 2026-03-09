package di

import (
	"github.com/michaeldcanady/go-onedrive/internal2/domain/account"
	domainauth "github.com/michaeldcanady/go-onedrive/internal2/domain/auth"
	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	domainenvironment "github.com/michaeldcanady/go-onedrive/internal2/domain/common/environment"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/config"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/editor"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	domainprofile "github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/state"
)

type Container interface {
	Cache() domaincache.Service2
	FS() domainfs.Service
	EnvironmentService() domainenvironment.EnvironmentService
	Logger() domainlogger.LoggerService
	Auth() domainauth.AuthService
	Profile() domainprofile.ProfileService
	Config() config.ConfigService
	State() state.Service
	Drive() drive.DriveService
	Account() account.Service
	Editor() editor.Service
}
