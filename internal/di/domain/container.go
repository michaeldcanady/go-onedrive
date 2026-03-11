package domain

import (
	domainaccount "github.com/michaeldcanady/go-onedrive/internal/account/domain"
	domainauth "github.com/michaeldcanady/go-onedrive/internal/auth/domain"
	domaincache "github.com/michaeldcanady/go-onedrive/internal/cache/domain"
	domainenv "github.com/michaeldcanady/go-onedrive/internal/core/env/domain"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	domainconfig "github.com/michaeldcanady/go-onedrive/internal/config/domain"
	domaindrive "github.com/michaeldcanady/go-onedrive/internal/drive/domain"
	domaineditor "github.com/michaeldcanady/go-onedrive/internal/editor/domain"
	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	domainprofile "github.com/michaeldcanady/go-onedrive/internal/profile/domain"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
)

type Container interface {
	Cache() domaincache.Service2
	FS() domainfs.Service
	IgnoreMatcherFactory() domainfs.IgnoreMatcherFactory
	EnvironmentService() domainenv.EnvironmentService
	Logger() domainlogger.LoggerService
	Auth() domainauth.AuthService
	Profile() domainprofile.ProfileService
	Config() domainconfig.ConfigService
	State() domainstate.Service
	Drive() domaindrive.DriveService
	Account() domainaccount.Service
	Editor() domaineditor.Service
}
