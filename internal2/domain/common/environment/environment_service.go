package environment

import "github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"

type EnvironmentService interface {
	CacheDir() (string, error)
	ConfigDir() (string, error)
	DataDir() (string, error)
	EnsureAll() error
	InstallDir() (string, error)
	IsLinux() bool
	IsMac() bool
	IsWindows() bool
	LogDir() (string, error)
	Name() string
	OS() string
	TempDir() (string, error)
	StateDir() (string, error)
	OutputDestination() (logger.OutputDestination, error)
	LogLevel() (string, error)
	Shell() (string, error)
	Editor() (string, error)
	Visual() (string, error)
}
