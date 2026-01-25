package di

import "context"

type EnvironmentService interface {
	IsWindows() bool
	IsMac() bool
	IsLinux() bool
	ConfigDir(context.Context) (string, error)
	DataDir(context.Context) (string, error)
	CacheDir(context.Context) (string, error)
	LogDir(context.Context) (string, error)
	InstallDir(context.Context) (string, error)
	TempDir(context.Context) (string, error)
	EnsureAll(context.Context) error
	Name() string
}
